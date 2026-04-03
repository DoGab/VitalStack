# Research: Food Datasources & Product Search Integration

> **GitHub Issue:** [#9 — Explore datasource for products](https://github.com/DoGab/VitalStack/issues/9)

## 🎯 Objective

Enable VitalStack users to log food not just by photo scanning, but also by **barcode scanning** and **free-text product search**. This requires integrating one or more external food/nutrition datasources that cover branded retail products (Migros, Coop, Lidl), product brands (Sponser, Emmi), and raw/generic foods (carrot, broccoli).

## 📋 Scope

- **Backend** — New service layer, external API clients, search index, new REST endpoints
- **Infrastructure** — Meilisearch (new component for product search, persistent on disk)
- **Frontend** — (Out of scope for this research — will be addressed in separate feature)
- **Frontend barcode scanner** — Separate issue (camera API + BarcodeDetector)

## ✅ Derived Definition of Done

- [ ] Backend can look up a product by EAN/barcode and return standardized nutrition data
- [ ] Backend can search products by free-text query and return ranked results
- [ ] Multiple datasources are queried in a waterfall/fan-out pattern for maximum coverage
- [ ] Results are normalized into a common `Product` domain type with macros
- [ ] Search is fast (<100ms for user-facing queries) and supports typo tolerance
- [ ] All datasources are free, self-hostable, and (preferably) open source
- [ ] Architecture follows the existing layered pattern (Controller → Service → Repository/Client)
- [ ] ODbL attribution for Open Food Facts data is handled

---

## 🔬 Datasource Comparison

### Candidate Datasources

| Datasource | Type | Coverage | Barcode? | Free-Text Search? | License | Swiss Coverage |
|---|---|---|---|---|---|---|
| **Open Food Facts (OFF)** | Crowdsourced | 4M+ global products | ✅ Native | ✅ API search | ODbL | ✅ Good (Migros, Coop present) |
| **USDA FoodData Central (FDC)** | Government | 450K+ (US focus) | ⚠️ Via UPC in Branded Foods | ✅ API search | Public Domain (CC0) | ❌ US-centric |
| **Swiss FSVO (naehrwertdaten.ch)** | Government | ~1,100 generic foods | ❌ No barcodes | ✅ REST API | Free (attribution) | ✅ Swiss reference data |
| **FoodRepo (EPFL)** | Academic/Crowdsourced | ~30K Swiss products | ✅ Barcode lookup | ✅ ElasticSearch API | CC BY 4.0 | ✅ Swiss-focused |

### Detailed Analysis

#### 1. Open Food Facts (OFF) — ⭐ PRIMARY DATASOURCE

- **URL:** https://world.openfoodfacts.org
- **API Base:** `https://world.openfoodfacts.org/api/v2/`
- **Go SDK:** [`github.com/openfoodfacts/openfoodfacts-go`](https://github.com/openfoodfacts/openfoodfacts-go)
- **Strengths:**
  - Largest open food database globally (4M+ products)
  - Native barcode lookup (`/api/v2/product/{barcode}`)
  - Good Swiss retailer coverage (Migros, Coop, Aldi, Lidl products)
  - Rich metadata: Nutri-Score, Eco-Score, NOVA groups, allergens, ingredients
  - Active community keeps data fresh
- **Weaknesses:**
  - Crowdsourced → variable quality/completeness per product
  - Some Swiss products may be missing or incomplete
  - Rate limits on public API (should respect usage guidelines)
- **License:** Open Database License (ODbL) — **must attribute source**
- **Self-hosting option:** Full MongoDB data dumps available for local import

#### 2. USDA FoodData Central (FDC) — ⭐ SECONDARY DATASOURCE (generic foods)

- **URL:** https://fdc.nal.usda.gov
- **API Base:** `https://api.nal.usda.gov/fdc/v1/`
- **Go SDK:** None official — use `net/http` + `encoding/json`
- **Strengths:**
  - Government-verified, laboratory-analyzed nutrient profiles
  - Excellent for generic/raw foods (Foundation Foods, SR Legacy datasets)
  - 190+ nutrients per food item (far beyond basic macros)
  - No rate limit concerns for reasonable usage (1,000/hr)
  - Public Domain (CC0) — no attribution required
- **Weaknesses:**
  - US-centric — no European/Swiss branded products
  - No native barcode scanning interface
  - Requires API key (free from data.gov)
- **Best for:** Raw foods (carrot, broccoli, chicken breast) where precision matters

#### 3. Swiss FSVO / naehrwertdaten.ch — 🔍 ENRICHMENT SOURCE

- **URL:** https://www.naehrwertdaten.ch
- **API Base:** `https://api.webapp.prod.blv.foodcase-services.com/BLV_WebApp_WS`
- **Go SDK:** None — custom HTTP client
- **Strengths:**
  - Official Swiss reference data for generic foods
  - Multilingual (DE/FR/IT/EN) — useful for Swiss users
  - Free REST API with no auth required
- **Weaknesses:**
  - Very small dataset (~1,100 generic foods only)
  - No barcodes / no branded products
  - Limited API documentation
- **Best for:** Swiss-specific names and reference values for common foods

#### 4. FoodRepo (EPFL) — 🔍 SWISS ENRICHMENT SOURCE

- **URL:** https://www.foodrepo.org
- **API Base:** `https://www.foodrepo.org/api/v3/`
- **Go SDK:** None — custom HTTP client
- **Strengths:**
  - Swiss-focused product database (~30K products)
  - Barcode lookup support
  - ElasticSearch-backed advanced search
  - Academic project with nutritional data from Swiss market
- **Weaknesses:**
  - Much smaller than OFF (~30K vs 4M)
  - Requires API key
  - Most Swiss products already covered by OFF
  - Unclear long-term maintenance commitment
- **Best for:** Supplementary Swiss-specific coverage if OFF has gaps

---

## 📊 Recommendation

### Multi-Source Strategy (Waterfall Pattern)

Use a **tiered datasource strategy** where we query sources in priority order:

```
User Request (Barcode or Search Query)
         │
         ▼
  ┌──────────────┐
  │  1. LOCAL     │  ← Meilisearch index (cached products)
  │     CACHE     │     Sub-millisecond response
  └──────┬───────┘
         │ miss
         ▼
  ┌──────────────┐
  │  2. OPEN      │  ← Primary external source
  │  FOOD FACTS   │     Barcode + text search
  └──────┬───────┘
         │ miss or generic food
         ▼
  ┌──────────────┐
  │  3. USDA      │  ← Best for raw/generic foods
  │     FDC       │     "carrot", "chicken breast"
  └──────┬───────┘
         │ miss (optional)
         ▼
  ┌──────────────┐
  │  4. SWISS     │  ← Swiss reference enrichment
  │    FSVO       │     (future / optional)
  └──────────────┘
```

**Why this order:**
1. **Local cache first** — instant response, no external API call
2. **OFF second** — widest coverage for branded products + barcodes
3. **USDA third** — superior accuracy for generic/raw foods
4. **FSVO optional** — Swiss-specific enrichment (Phase 2)

> **FoodRepo (EPFL) is excluded** from the primary strategy because its coverage largely overlaps with OFF but at a much smaller scale. If Swiss-specific gaps emerge in OFF, FoodRepo can be added as an enrichment layer later.

---

## 🏗️ Architecture Proposal

### New Components

#### 1. Meilisearch — Product Search Engine (NEW INFRASTRUCTURE)

- **Purpose:** Fast, typo-tolerant product search for user-facing queries
- **Why not Postgres FTS?** Product search is a core feature — users expect instant, typo-tolerant "search-as-you-type" results. Meilisearch provides this out-of-the-box.
- **Self-hosted:** Docker container, single binary, MIT license
- **Resource footprint:** ~100MB RAM for 100K products
- **Go SDK:** [`github.com/meilisearch/meilisearch-go`](https://github.com/meilisearch/meilisearch-go)

Docker addition to `docker-compose.yml`:
```yaml
services:
  meilisearch:
    image: getmeili/meilisearch:v1.12
    container_name: vitalstack-meilisearch
    restart: unless-stopped
    ports:
      - "7700:7700"
    volumes:
      - meilisearch_data:/meili_data
    environment:
      - MEILI_MASTER_KEY=${MEILI_MASTER_KEY}
      - MEILI_ENV=production
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:7700/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  meilisearch_data:
```

#### 2. Normalized Product Domain Model

A unified `Product` type that all datasources normalize into:

```go
// Product represents a food product from any datasource
type Product struct {
    ID          string   `json:"id"`          // Internal ID or barcode
    Barcode     string   `json:"barcode"`     // EAN/UPC barcode
    Name        string   `json:"name"`
    Brand       string   `json:"brand"`
    Categories  []string `json:"categories"`
    ImageURL    string   `json:"image_url"`
    Source      string   `json:"source"`      // "openfoodfacts", "usda", "fsvo"
    NutriScore  string   `json:"nutri_score"` // A-E (OFF only)
    Macros      MacrosPer100g `json:"macros"`
}

type MacrosPer100g struct {
    Calories float64 `json:"calories"`
    Protein  float64 `json:"protein"`
    Carbs    float64 `json:"carbs"`
    Fat      float64 `json:"fat"`
    Fiber    float64 `json:"fiber"`
}
```

### Proposed Backend Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        HTTP Request                              │
│        GET /api/products/search?q=yogurt                         │
│        GET /api/products/barcode/{ean}                           │
└───────────────────────────┬─────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                   ProductController (NEW)                         │
│              internal/controller/product_controller.go            │
│   • Registers Huma routes for search + barcode lookup            │
│   • Maps HTTP DTOs ↔ Service types                               │
└───────────────────────────┬─────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                   ProductService (NEW)                            │
│                  pkg/service/product_service.go                   │
│   • Orchestrates multi-source lookup                             │
│   • Implements waterfall pattern (cache → OFF → USDA)            │
│   • Normalizes results into common Product type                  │
│   • Populates Meilisearch cache on successful external lookups   │
└───────────────┬──────────────────────┬──────────────────────────┘
                │                      │
      ┌─────────▼──────────┐  ┌───────▼─────────────────┐
      │  Search Index        │  │  External API Clients   │
      │  (Meilisearch)       │  │                         │
      │  pkg/search/         │  │  pkg/datasource/        │
      │  search_client.go    │  │  ├── off_client.go      │
      │                      │  │  ├── usda_client.go     │
      │  • Index products    │  │  └── client.go (iface) │
      │  • Search products   │  │                         │
      │  • Typo-tolerant     │  │  • HTTP clients for     │
      │    full-text search  │  │    external food APIs    │
      └─────────────────────┘  │  • Normalize to Product  │
                                └─────────────────────────┘
```

### New Files to Create

```
apps/api-go/
├── internal/
│   └── controller/
│       ├── product_controller.go   [NEW] HTTP handlers for search & barcode
│       └── product_types.go        [NEW] Request/Response DTOs
├── pkg/
│   ├── service/
│   │   └── product_service.go      [NEW] Business logic, waterfall orchestration
│   ├── datasource/                 [NEW] Package for external API clients
│   │   ├── client.go               [NEW] FoodDatasource interface
│   │   ├── off_client.go           [NEW] Open Food Facts API client
│   │   └── usda_client.go          [NEW] USDA FoodData Central API client
│   ├── search/                     [NEW] Package for Meilisearch integration
│   │   └── meilisearch_client.go   [NEW] Meilisearch index/search operations
│   └── types/
│       └── product.go              [NEW] Domain Product type
```

### Files to Modify

```
apps/api-go/
├── cmd/server.go                   [MODIFY] Wire up new ProductController + dependencies
├── internal/conf/conf.go           [MODIFY] Add Meilisearch + USDA API key config flags
├── local-config.yaml               [MODIFY] Add local Meilisearch + USDA config

docker-compose.yml                  [MODIFY] Add Meilisearch service
docker-compose.local.yml            [MODIFY] Add local Meilisearch override
```

### New API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/products/barcode/{ean}` | Look up product by EAN/UPC barcode |
| `GET` | `/api/products/search?q={query}&limit={n}` | Search products by name/brand |

---

## 🗺️ Discovered Codebase State

### Existing Patterns

- `apps/api-go/internal/controller/nutrition_controller.go:33-78` — Controller pattern: registers Huma operations for all routes in `Register(api huma.API)`. Each handler converts between controller DTOs and service types.
- `apps/api-go/pkg/service/nutrition_service.go:1-58` — Service pattern: struct with injected dependencies (Genkit, FoodLogRepository), functional options pattern (`WithMockScan`), constructor `NewNutritionService()`.
- `apps/api-go/internal/controller/nutrition_types.go:1-128` — DTO pattern: Input structs embed `Body` pointer for Huma, Output structs have `Body` pointer, all fields decorated with `json`, `doc`, `example` tags.
- `apps/api-go/cmd/server.go:24-90` — Wiring pattern: dependencies created in `ServerEntryPoint()`, injected into services, then into controllers, registered via `api.RegisterAPI(ctrl)`.
- `apps/api-go/internal/conf/conf.go:1-136` — Config pattern: consts for flag names/defaults/help, registered via Cobra PersistentFlags, bound to Viper.
- `apps/api-go/internal/repository/food_log.go:1-175` — Repository pattern: interface + struct implementation, Supabase client as dependency.

### Integration Points

- `apps/api-go/cmd/server.go:49-67` — Where new `ProductController` gets instantiated and registered with the API server
- `apps/api-go/internal/conf/conf.go:102-107` — Where new config constants for Meilisearch URL/key and USDA API key will be added
- `apps/api-go/internal/server/server.go:134-138` — `RegisterAPI()` accepts variadic `Controller` args — new controller plugs in here
- `docker-compose.yml:7-58` — Where Meilisearch service gets added
- `apps/api-go/local-config.yaml` — Where local dev Meilisearch/USDA config gets added

### Key Interfaces to Implement

```go
// Follows existing Controller interface: internal/server/server.go:26
type Controller interface {
    Register(api huma.API)
}

// New interface for datasource abstraction
type FoodDatasource interface {
    LookupBarcode(ctx context.Context, barcode string) (*Product, error)
    SearchProducts(ctx context.Context, query string, limit int) ([]Product, error)
}

// New interface for search index
type ProductSearchIndex interface {
    IndexProduct(ctx context.Context, product Product) error
    IndexProducts(ctx context.Context, products []Product) error
    Search(ctx context.Context, query string, limit int) ([]Product, error)
}
```

---

## 🔗 Architecture References

- Backend: `apps/api-go/architecture.md`
- Frontend: `apps/web/architecture.md`
- Server wiring: `apps/api-go/cmd/server.go`
- Config system: `apps/api-go/internal/conf/conf.go`
- Docker infra: `docker-compose.yml`, `docker-compose.local.yml`

---

## ✅ Resolved Decisions

1. **Meilisearch data seeding:** ✅ **Pre-seed** the Meilisearch index with an Open Food Facts data dump (JSONL format). This ensures search works immediately without waiting for lazy population. A seeding script/command will be needed.
2. **Barcode scanner:** ✅ **Separate issue.** The frontend barcode scanner (camera API + BarcodeDetector) is out of scope for this feature. The backend API endpoints (`/api/products/barcode/{ean}`) will be ready for when that is implemented.
3. **Product caching in Supabase:** ✅ **No Supabase cache needed.** Meilisearch is persistent — it uses memory-mapped files (LMDB) on disk. Data survives container/process restarts as long as the Docker volume (`meilisearch_data:/meili_data`) is mounted, which is already configured. Meilisearch serves as both the search engine AND the product store.
4. **USDA API key management:** ✅ **Viper config.** Store the USDA API key in the Viper configuration system alongside the Gemini key. Add `usda.api-key` config flag with env var `USDA_API_KEY`.
5. **Data ownership boundary (Supabase vs Meilisearch) — Copy vs Reference:**

   | Data | Store | Reason |
   |------|-------|--------|
   | User food logs (`food_logs`, `food_log_ingredients`) | **Supabase (Postgres)** | User-scoped, transactional, RLS-protected, relational joins, time-range aggregation |
   | User profiles (`profiles`) | **Supabase (Postgres)** | Auth-linked, RLS-protected |
   | Product catalog (OFF/USDA imports) | **Meilisearch** | Shared reference data, read-only, typo-tolerant full-text search, no user ownership |

   **Should food_log_ingredients reference the product or copy the macro data?**

   Three options were evaluated:

   | Approach | How it works | Verdict |
   |----------|-------------|---------|
   | **A) Pure reference** — store only `product_id`, fetch macros from Meilisearch at read time | `food_log_ingredients.product_id` → Meilisearch lookup | ❌ **Not viable** |
   | **B) Pure copy (current pattern)** — snapshot macros into `food_log_ingredients` at log time | Macros stored inline, no product reference | ✅ Works, but no link back to product |
   | **C) Hybrid** — snapshot macros AND store an optional `product_id` for linking | Both inline macros + optional `source_product_id` | ✅ **Recommended** |

   **Why pure referencing (Option A) doesn't work:**
   1. **No cross-system JOINs** — Postgres and Meilisearch are separate systems. To render the daily intake page, you'd need: query Postgres for logs → for each ingredient, query Meilisearch by product ID → assemble in Go. That's N+1 network calls across two systems per page load. This is *less* scalable than denormalization.
   2. **Historical accuracy breaks** — If OFF updates "Emmi Caffè Latte" from 180→175 kcal, past logs would silently change. A food log is an *event record*: "I ate X with Y macros at time Z." That snapshot must be immutable.
   3. **Meilisearch downtime = broken food history** — If Meilisearch is being re-seeded or restarted, the entire food log UI breaks because it can't resolve product references.
   4. **Re-seeding destroys references** — When you re-import an OFF dump, product IDs are not guaranteed stable. All FK references become orphaned.
   5. **Storage is negligible** — A `food_log_ingredient` row is ~200 bytes. 10 items/day × 365 days × 10 years = ~7MB per user. Insignificant.

   **Why hybrid (Option C) is recommended:**
   - ✅ Keeps the immutable macro snapshot for historical accuracy (same as today's AI scan flow, see `nutrition_service.go:186-198`)
   - ✅ Adds an **optional** `source_product_id TEXT` column — non-breaking, nullable
   - ✅ Enables future "log again" / "view original product" UX flows
   - ✅ Allows analytics (most-logged products) without breaking if Meilisearch is rebuilt
   - ✅ Meilisearch can be re-seeded without affecting any user data — stores are fully decoupled
   - ✅ Daily intake page is a single Postgres query with no external dependencies

---

> **Research complete. Ready for Planning phase.**
