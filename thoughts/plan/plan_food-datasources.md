# Plan: Food Datasources & Product Search Integration

> Source: `thoughts/research/research_food-datasources.md`

## 🎯 Objective

Enable VitalStack users to look up food products by **barcode** and **free-text search** by integrating Open Food Facts (primary) and USDA FoodData Central (secondary) as external datasources, with Meilisearch providing fast, typo-tolerant search over a cached product index.

## 🚫 What We're NOT Doing

- **Frontend UI** for barcode scanning or product search (separate feature)
- **Frontend barcode scanner** (camera API / BarcodeDetector — separate issue)
- **Swiss FSVO / FoodRepo integration** (deferred enrichment source)
- **Bulk OFF data seeding** — products are lazily cached into Meilisearch on first lookup; a seeding script is a follow-up task
- **Product image storage/proxy** — we store `image_url` only, no image caching
- **Supabase product table** — Meilisearch is the sole product store (per research decision #5)
- **`source_product_id` column migration** — the hybrid approach from research is deferred; current `food_log_ingredients` schema remains unchanged. Will be added when the product search UI feature is built.

## ✅ Definition of Done

- [ ] Backend can look up a product by EAN/barcode and return standardized nutrition data
- [ ] Backend can search products by free-text query and return ranked results
- [ ] Multiple datasources are queried in a waterfall pattern for maximum coverage
- [ ] Results are normalized into a common `Product` domain type with macros per 100g
- [ ] Search is fast (<100ms for user-facing queries) and supports typo tolerance
- [ ] All datasources are free, self-hostable, and (preferably) open source
- [ ] Architecture follows the existing layered pattern (Controller → Service → Client)
- [ ] ODbL attribution for Open Food Facts data is included in responses
- [ ] Every new exported function/type has a Go doc comment
- [ ] Every new package has unit tests

---

## 📝 Implementation Phases

### Phase 1: Infrastructure — Meilisearch + Configuration

**Overview:** Add Meilisearch to the Docker infrastructure and wire its connection config (plus USDA API key) into the Go backend's Viper-based configuration system. Zero application logic — only infrastructure foundation.

**Changes Required:**

- `docker-compose.yml` — Add `meilisearch` service block:
  - Image: `getmeili/meilisearch:v1.12`
  - Port: `7700:7700`
  - Volume: `meilisearch_data:/meili_data`
  - Environment: `MEILI_MASTER_KEY=${MEILI_MASTER_KEY}`, `MEILI_ENV=production`
  - Health check: `curl -f http://localhost:7700/health`
  - Add `meilisearch_data` to the `volumes:` section

- `docker-compose.local.yml` — Add Meilisearch local override:
  - `MEILI_ENV=development`
  - No master key required for local dev (`MEILI_MASTER_KEY=` empty or omit)

- `apps/api-go/internal/conf/conf.go` — Add new config constants following the existing pattern (const group + `RegisterFlags`):
  ```
  meilisearch.url       → default "http://localhost:7700"
  meilisearch.api-key   → default ""
  usda.api-key          → default ""
  ```
  Register all three as `PersistentFlags` in `RegisterFlags()` and bind to Viper.

- `apps/api-go/local-config.yaml` — Add:
  ```yaml
  meilisearch:
    url: "http://localhost:7700"
    api-key: ""

  usda:
    api-key: ""  # Get free key from https://api.data.gov/signup/
  ```

**Success Criteria:**

#### Automated Verification:
- [x] `cd apps/api-go && go build ./...` compiles without errors
- [x] `docker compose -f docker-compose.yml config` validates without errors
- [x] `cd apps/api-go && golangci-lint run ./...` passes

#### Manual Verification:
- [x] `docker compose up meilisearch -d` starts successfully
- [x] `curl http://localhost:7700/health` returns `{"status":"available"}`

> **Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation that Meilisearch is reachable before proceeding.

---

### Phase 2: Domain Types + External API Clients

**Overview:** Create the normalized `Product` domain model and build HTTP clients for Open Food Facts and USDA FoodData Central. Each client implements a shared `FoodDatasource` interface and normalizes external API responses into the common `Product` type. This is pure library code with unit tests — no server wiring.

**Changes Required:**

#### 2a. Domain types

- `apps/api-go/pkg/types/product.go` — **[NEW]** Domain types:
  ```go
  // Product represents a normalized food product from any datasource.
  type Product struct {
      ID         string        `json:"id"`
      Barcode    string        `json:"barcode"`
      Name       string        `json:"name"`
      Brand      string        `json:"brand"`
      Categories []string      `json:"categories"`
      ImageURL   string        `json:"image_url"`
      Source     string        `json:"source"`      // "openfoodfacts", "usda"
      NutriScore string        `json:"nutri_score"` // A-E (OFF only)
      Macros     MacrosPer100g `json:"macros"`
  }

  // MacrosPer100g holds nutritional data normalized to per-100g values.
  type MacrosPer100g struct {
      Calories float64 `json:"calories"`
      Protein  float64 `json:"protein"`
      Carbs    float64 `json:"carbs"`
      Fat      float64 `json:"fat"`
      Fiber    float64 `json:"fiber"`
  }
  ```
  All fields with `json` tags. `doc` and `example` tags are NOT needed here — these are domain types, not Huma DTOs.

#### 2b. Datasource interface

- `apps/api-go/pkg/datasource/client.go` — **[NEW]** Shared interface:
  ```go
  // FoodDatasource defines the contract for any external food data provider.
  type FoodDatasource interface {
      LookupBarcode(ctx context.Context, barcode string) (*types.Product, error)
      SearchProducts(ctx context.Context, query string, limit int) ([]types.Product, error)
      Name() string // Returns datasource name e.g. "openfoodfacts"
  }
  ```
  Define a sentinel `ErrNotFound` error for barcode misses (allows the waterfall to continue without logging errors).

#### 2c. Open Food Facts client

- `apps/api-go/pkg/datasource/off_client.go` — **[NEW]**:
  - Constructor: `NewOFFClient(httpClient *http.Client) *OFFClient`
  - Use `net/http` directly (the official Go SDK is thin/unmaintained — raw HTTP is more reliable)
  - Set `User-Agent: VitalStack/1.0 (github.com/DoGab/VitalStack)` per OFF guidelines
  - Barcode: `GET https://world.openfoodfacts.org/api/v2/product/{barcode}?fields=code,product_name,brands,categories_tags,image_url,nutriments,nutriscore_grade`
  - Search: `GET https://world.openfoodfacts.org/cgi/search.pl?search_terms={q}&json=1&page_size={limit}&fields=...`
  - Internal types for OFF JSON response (unexported) → normalize to `types.Product`

- `apps/api-go/pkg/datasource/off_client_test.go` — **[NEW]**:
  - Use `httptest.NewServer` to mock OFF API responses
  - Test fixtures: create `apps/api-go/pkg/datasource/testdata/off_product.json` and `testdata/off_search.json` with representative OFF response payloads
  - Test cases: successful barcode lookup, barcode not found (status 0), empty search results, malformed JSON, missing nutriment fields

#### 2d. USDA FoodData Central client

- `apps/api-go/pkg/datasource/usda_client.go` — **[NEW]**:
  - Constructor: `NewUSDAClient(httpClient *http.Client, apiKey string) *USDAClient`
  - Search: `GET https://api.nal.usda.gov/fdc/v1/foods/search?query={q}&pageSize={limit}&api_key={key}`
  - Barcode: `GET https://api.nal.usda.gov/fdc/v1/foods/search?query={upc}&dataType=Branded&api_key={key}` — filter results by matching `gtinUpc` field
  - Internal types for FDC JSON response (unexported) → normalize to `types.Product`
  - Nutrient extraction: iterate `foodNutrients` array, match by `nutrientId` (energy=1008, protein=1003, carbs=1005, fat=1004, fiber=1079)

- `apps/api-go/pkg/datasource/usda_client_test.go` — **[NEW]**:
  - Use `httptest.NewServer` to mock USDA API
  - Test fixtures: `testdata/usda_search.json`
  - Test cases: successful search, API key injection, empty results, nutrient ID mapping

**Success Criteria:**

#### Automated Verification:
- [x] `cd apps/api-go && go test ./pkg/types/... ./pkg/datasource/...` — all tests pass
- [x] `cd apps/api-go && golangci-lint run ./...` passes
- [x] `cd apps/api-go && go vet ./...` passes

#### Manual Verification:
- [ ] N/A — pure library code

> **Implementation Note**: No pause needed — proceed directly to Phase 3.

---

### Phase 3: Meilisearch Search Index Client

**Overview:** Build the Meilisearch integration layer that indexes and searches products. Wraps the official `meilisearch-go` SDK behind a `ProductSearchIndex` interface for testability. Includes adding the Go dependency.

**Changes Required:**

- Run `cd apps/api-go && go get github.com/meilisearch/meilisearch-go` — add SDK dependency

- `apps/api-go/pkg/search/search_client.go` — **[NEW]**:
  - `ProductSearchIndex` interface:
    ```go
    // ProductSearchIndex defines search and indexing operations for the product catalog.
    type ProductSearchIndex interface {
        IndexProduct(ctx context.Context, product types.Product) error
        IndexProducts(ctx context.Context, products []types.Product) error
        Search(ctx context.Context, query string, limit int) ([]types.Product, error)
        LookupBarcode(ctx context.Context, barcode string) (*types.Product, error)
    }
    ```
  - `MeilisearchClient` struct implementing the interface
  - Constructor: `NewMeilisearchClient(url, apiKey string) (*MeilisearchClient, error)`
  - On init: create or get `products` index, configure searchable attributes (`name`, `brand`, `categories`), filterable attributes (`barcode`, `source`), sortable attributes (`name`)
  - Primary key strategy: use `id` field. For OFF products, `id = "off-{barcode}"`. For USDA, `id = "usda-{fdcId}"`. This prevents collisions and allows deduplication.
  - `Search`: calls `Index("products").Search(query, &SearchRequest{Limit: limit})`
  - `LookupBarcode`: calls `Index("products").Search("", &SearchRequest{Filter: "barcode = '{ean}'", Limit: 1})`
  - `IndexProduct`: upserts single product via `AddDocuments`
  - `IndexProducts`: bulk upserts via `AddDocuments`

- `apps/api-go/pkg/search/search_client_test.go` — **[NEW]**:
  - Unit tests for result mapping and configuration logic
  - For integration tests (optional): test against running Meilisearch with build tag `//go:build integration`

**Success Criteria:**

#### Automated Verification:
- [x] `cd apps/api-go && go build ./...` compiles
- [x] `cd apps/api-go && go test ./pkg/search/...` — unit tests pass
- [x] `cd apps/api-go && golangci-lint run ./...` passes

#### Manual Verification:
- [ ] N/A — integration tested in Phase 4

> **Implementation Note**: No pause needed — proceed directly to Phase 4.

---

### Phase 4: ProductService + ProductController + Server Wiring

**Overview:** Create the `ProductService` (waterfall orchestrator), the `ProductController` (Huma HTTP endpoints), and wire everything into the server's dependency injection in `cmd/server.go`. This makes the API endpoints live.

**Changes Required:**

#### 4a. ProductService

- `apps/api-go/pkg/service/product_service.go` — **[NEW]**:
  - `ProductService` struct with injected deps: `search.ProductSearchIndex`, `[]datasource.FoodDatasource`
  - Constructor: `NewProductService(index search.ProductSearchIndex, datasources ...datasource.FoodDatasource) *ProductService`
  - `LookupBarcode(ctx context.Context, barcode string) (*types.Product, error)`:
    1. Query search index by barcode filter
    2. On miss → iterate datasources in order (OFF first, USDA second)
    3. On external hit → index the product in Meilisearch (fire-and-forget goroutine)
    4. If all sources miss → return `ErrNotFound`
  - `SearchProducts(ctx context.Context, query string, limit int) ([]types.Product, error)`:
    1. Search Meilisearch index
    2. If results < limit → fan out to external datasources, merge, deduplicate by barcode
    3. Cache new products found externally into Meilisearch
    4. Return merged results, capped at `limit`

- `apps/api-go/pkg/service/product_service_test.go` — **[NEW]**:
  - Mock `ProductSearchIndex` and `FoodDatasource` via interfaces
  - Test: index hit returns immediately without calling external sources
  - Test: index miss cascades to OFF, then USDA
  - Test: external hit triggers indexing callback
  - Test: deduplication when same barcode found in both index and external
  - Test: all sources miss → returns appropriate error

#### 4b. ProductController + DTOs

- `apps/api-go/internal/controller/product_types.go` — **[NEW]** Request/Response DTOs:
  ```go
  // BarcodeInput represents a barcode lookup request
  type BarcodeInput struct {
      EAN string `path:"ean" doc:"EAN/UPC barcode" example:"7613035466432"`
  }

  // SearchProductsInput represents a product search request
  type SearchProductsInput struct {
      Query string `query:"q" required:"true" doc:"Search query" example:"yogurt"`
      Limit int    `query:"limit" default:"10" doc:"Max results to return" example:"10"`
  }

  // ProductBody represents a product in API responses
  type ProductBody struct {
      ID         string            `json:"id" doc:"Product identifier" example:"off-7613035466432"`
      Barcode    string            `json:"barcode" doc:"EAN/UPC barcode" example:"7613035466432"`
      Name       string            `json:"name" doc:"Product name" example:"Caffè Latte"`
      Brand      string            `json:"brand" doc:"Brand name" example:"Emmi"`
      ImageURL   string            `json:"image_url,omitempty" doc:"Product image URL"`
      Source     string            `json:"source" doc:"Data source" example:"openfoodfacts"`
      NutriScore string            `json:"nutri_score,omitempty" doc:"Nutri-Score grade (A-E)" example:"C"`
      Macros     MacrosPer100gBody `json:"macros" doc:"Nutritional values per 100g"`
  }

  // MacrosPer100gBody represents per-100g nutritional data in API responses
  type MacrosPer100gBody struct {
      Calories float64 `json:"calories" doc:"Calories per 100g" example:"56"`
      Protein  float64 `json:"protein" doc:"Protein per 100g (grams)" example:"3.2"`
      Carbs    float64 `json:"carbs" doc:"Carbohydrates per 100g (grams)" example:"5.1"`
      Fat      float64 `json:"fat" doc:"Fat per 100g (grams)" example:"1.8"`
      Fiber    float64 `json:"fiber" doc:"Fiber per 100g (grams)" example:"0"`
  }

  // BarcodeOutput is the HTTP response for barcode lookups
  type BarcodeOutput struct {
      Body *ProductBody
  }

  // SearchProductsOutput is the HTTP response for product search
  type SearchProductsOutput struct {
      Body *SearchProductsOutputBody
  }

  // SearchProductsOutputBody wraps the list of products with attribution
  type SearchProductsOutputBody struct {
      Products    []ProductBody `json:"products" doc:"List of matching products"`
      Attribution string        `json:"attribution,omitempty" doc:"Data attribution notice"`
  }
  ```
  Include conversion helper: `productBodyFromDomain(p types.Product) ProductBody`

- `apps/api-go/internal/controller/product_controller.go` — **[NEW]**:
  - Define `ProductServicer` interface (for testability):
    ```go
    type ProductServicer interface {
        LookupBarcode(ctx context.Context, barcode string) (*types.Product, error)
        SearchProducts(ctx context.Context, query string, limit int) ([]types.Product, error)
    }
    ```
  - `ProductController` struct with `ProductServicer` field
  - Constructor: `NewProductController(svc ProductServicer) *ProductController`
  - `Register(api huma.API)` — register two operations:
    - `GET /api/products/barcode/{ean}` → `BarcodeHandler`, tag `products`, operationID `lookup-product-barcode`
    - `GET /api/products/search` → `SearchHandler`, tag `products`, operationID `search-products`
  - `BarcodeHandler`: call service, convert domain → DTO, return. On `ErrNotFound` → `huma.Error404NotFound`
  - `SearchHandler`: call service, convert domain → DTOs, include ODbL `attribution` string if any OFF results present

#### 4c. Server wiring

- `apps/api-go/cmd/server.go` — **[MODIFY]** Add wiring after the existing nutrition controller block (around line 64):
  1. Import: `"github.com/dogab/vitalstack/api/pkg/datasource"`, `"github.com/dogab/vitalstack/api/pkg/search"`
  2. Create Meilisearch client:
     ```go
     meiliClient, err := search.NewMeilisearchClient(
         viper.GetString(conf.MeilisearchURLArg),
         viper.GetString(conf.MeilisearchAPIKeyArg),
     )
     ```
     Handle error with `slog.Warn` (non-fatal — product search is optional)
  3. Create datasource clients:
     ```go
     offClient := datasource.NewOFFClient(http.DefaultClient)
     usdaClient := datasource.NewUSDAClient(http.DefaultClient, viper.GetString(conf.USDAAPIKeyArg))
     ```
  4. Create service and controller:
     ```go
     productSvc := service.NewProductService(meiliClient, offClient, usdaClient)
     productCtrl := controller.NewProductController(productSvc)
     ```
  5. Register: `api.RegisterAPI(ctrl, productCtrl)`

**Success Criteria:**

#### Automated Verification:
- [x] `cd apps/api-go && go test ./...` — all tests pass (existing + new)
- [x] `cd apps/api-go && golangci-lint run ./...` passes
- [x] `cd apps/api-go && go build ./...` compiles
- [x] `make openapi` succeeds — OpenAPI spec includes new product endpoints

#### Manual Verification:
- [ ] Start Meilisearch: `docker compose up meilisearch -d`
- [ ] Start backend: `make dev-api` (or `make dev` for full stack)
- [ ] `curl http://localhost:8080/api/products/barcode/7613035466432` returns a product with macros (first call hits OFF, subsequent calls are cached)
- [ ] `curl "http://localhost:8080/api/products/search?q=yogurt&limit=5"` returns search results
- [ ] API docs at `http://localhost:8080/docs` show the new `products` tag with both endpoints
- [ ] Second barcode lookup for same EAN is noticeably faster (served from Meilisearch)

> **Implementation Note**: After completing this phase and all automated verification passes, **pause here for manual confirmation** from the human that the endpoints work correctly before proceeding to cleanup.

---

### Phase 5: Final Cleanup & Documentation

**Overview:** Regenerate the TypeScript client, update architecture documentation, and do a final verification pass.

**Changes Required:**

- Run `make openapi` — regenerate OpenAPI YAML spec + TypeScript client in `apps/web`
- Verify TypeScript client in `apps/web/src/lib/api/` reflects new `lookupProductBarcode` and `searchProducts` methods

- `apps/api-go/architecture.md` — **[MODIFY]** Add:
  - `pkg/datasource/` package description (FoodDatasource interface, OFF + USDA clients)
  - `pkg/search/` package description (ProductSearchIndex, Meilisearch wrapper)
  - Meilisearch as an infrastructure component in the architecture diagram
  - Note about waterfall lookup pattern

- `README.md` or `.env.example` — **[MODIFY]** Document new environment variables:
  - `MEILI_MASTER_KEY` — Meilisearch master key (production only)
  - `USDA_API_KEY` — free key from https://api.data.gov/signup/

**Success Criteria:**

#### Automated Verification:
- [ ] `make openapi` runs without errors
- [ ] TypeScript client in `apps/web` has generated types for product endpoints
- [ ] `cd apps/api-go && go test ./...` — all tests still pass
- [ ] `cd apps/api-go && golangci-lint run ./...` passes

#### Manual Verification:
- [ ] Architecture docs accurately describe the new packages and data flow
- [ ] `docker compose -f docker-compose.yml config` validates with Meilisearch service
