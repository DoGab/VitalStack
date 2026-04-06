# Issue #8: Barcode Scanning & Product Search — Implementation Plan

## Overview

Add barcode scanning and product search flows to VitalStack's frontend, with per-request language propagation through the backend waterfall. Users can search for products by name or scan barcodes, view macro data, select serving sizes, and log products using the existing `/api/nutrition/log` endpoint.

## Current State Analysis

### Backend
- **Waterfall architecture**: `Meilisearch → OFF → FSVO → USDA` orchestrated by `ProductService`
- **Language is boot-time only**: `OFFClient.language` and `FSVOClient.language` are struct fields set once via functional options at startup
- **`FoodDatasource` interface**: `LookupBarcode(ctx, barcode)` and `SearchProducts(ctx, query, limit)` — no `lang` param
- **`ProductServicer` interface**: mirrors the same signature — `LookupBarcode(ctx, barcode)`, `SearchProducts(ctx, query, limit)`
- **`Product` type**: has no `ServingSize` or `ServingQuantity` fields — macros are per-100g only
- **OFF fields**: `offFields = "code,product_name,brands,categories_tags,image_url,nutriments,nutriscore_grade"` — does not request `serving_size` or `serving_quantity`
- **Log endpoint**: `POST /api/nutrition/log` requires `confidence` (float64, non-optional) and `ingredients` (non-optional array)

### Frontend
- **`AddEntryModal.svelte`**: currently has 3 options (camera, upload, quick add) — no search or barcode scan
- **`FoodScannerModal.svelte`**: handles camera → AI scan → results → log flow — pattern to follow for new modals
- **API client**: `openapi-fetch` with typed paths from auto-generated `schema.d.ts`
- **No barcode scanning library** installed

## Desired End State

1. User can search products by name in a `ProductSearchModal` with debounced API calls
2. User can scan barcodes via camera in a `BarcodeScannerModal` (chosen library from evaluation)
3. Product results show macros per 100g, brand, image, and source
4. User can pick a serving size (100g, product serving, or custom weight) in `ProductLogDrawer`
5. Product logs go through existing `/api/nutrition/log` endpoint and count toward daily totals
6. Backend accepts `lang` query parameter and propagates it per-request to OFF and FSVO
7. `Product` type includes `serving_size` (string) and `serving_quantity` (float64) extracted from OFF

## What We're NOT Doing

- No new `/api/products/log` endpoint — reusing existing `/api/nutrition/log` (Option B)
- No database schema changes — product logs are stored as food logs with the same schema
- No barcode scanning on desktop (mobile-only camera feature)
- No multi-language product name translation — `lang` only affects which locale OFF/FSVO use for responses
- No offline barcode scanning — requires network for product lookup after scan

## Design Decisions

1. **Barcode library**: Side-by-side evaluation in Phase 2 (`html5-qrcode` vs `@niczvt/barcode-scanner-svelte`)
2. **Product logging**: Adapt existing `/api/nutrition/log` — make `confidence` and `ingredients` optional so product-based logs (no AI scan) work alongside AI scan logs   
3. **Serving size**: Add `serving_size` (string, e.g. "250ml") and `serving_quantity` (float64, e.g. 250) to `Product` type, extract from OFF. FSVO/USDA default to empty/100g

---

## Phase 1: Backend — Per-request Language + Product Enhancements

### Overview
Thread a `lang` query parameter from the API endpoints through the service layer to the datasource clients. Add serving size fields to the `Product` type and extract them from OFF. Make the log endpoint accept product-based logs.

### Changes Required:

#### 1. FoodDatasource Interface — Add `lang` parameter
**File**: [`client.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/client.go)
**Changes**: Add `lang string` to `LookupBarcode` and `SearchProducts`.

```go
type FoodDatasource interface {
    Name() string
    LookupBarcode(ctx context.Context, barcode string, lang string) (*types.Product, error)
    SearchProducts(ctx context.Context, query string, limit int, lang string) ([]types.Product, error)
}
```

#### 2. OFF Client — Accept per-request language
**File**: [`off_client.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/off_client.go)
**Changes**:
- Update method signatures to accept `lang string`
- Use the per-request `lang` when non-empty, falling back to `c.language`
- Add `serving_size` and `serving_quantity` to `offFields` constant
- Add `ServingSize` and `ServingQuantity` fields to `offProduct` struct
- Map new fields in `offProductToDomain`

```go
const offFields = "code,product_name,brands,categories_tags,image_url,nutriments,nutriscore_grade,serving_size,serving_quantity"

type offProduct struct {
    // ... existing fields ...
    ServingSize     string  `json:"serving_size"`
    ServingQuantity float64 `json:"serving_quantity"`
}

func (c *OFFClient) LookupBarcode(ctx context.Context, barcode string, lang string) (*types.Product, error) {
    effectiveLang := c.language
    if lang != "" {
        effectiveLang = lang
    }
    // ... use effectiveLang instead of c.language ...
}

func (c *OFFClient) SearchProducts(ctx context.Context, query string, limit int, lang string) ([]types.Product, error) {
    effectiveLang := c.language
    if lang != "" {
        effectiveLang = lang
    }
    // ... use effectiveLang instead of c.language ...
}
```

#### 3. FSVO Client — Accept per-request language
**File**: [`fsvo_client.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/fsvo_client.go)
**Changes**: Same pattern — accept `lang string`, use when non-empty.

```go
func (c *FSVOClient) LookupBarcode(_ context.Context, _ string, _ string) (*types.Product, error) {
    return nil, ErrNotFound
}

func (c *FSVOClient) SearchProducts(ctx context.Context, query string, limit int, lang string) ([]types.Product, error) {
    effectiveLang := c.language
    if lang != "" {
        effectiveLang = lang
    }
    // ... pass effectiveLang to searchFoods and getFood ...
}
```

**Also update internal methods** `searchFoods` and `getFood` to accept `lang string` parameter instead of using `c.language`.

#### 4. USDA Client — Accept per-request language (ignored)
**File**: [`usda_client.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/usda_client.go)
**Changes**: Accept `lang string` parameter to satisfy interface, but ignore it (USDA has no language support).

```go
func (c *USDAClient) LookupBarcode(ctx context.Context, barcode string, _ string) (*types.Product, error) { ... }
func (c *USDAClient) SearchProducts(ctx context.Context, query string, limit int, _ string) ([]types.Product, error) { ... }
```

#### 5. Product Type — Add serving fields
**File**: [`product.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/types/product.go)
**Changes**: Add `ServingSize` and `ServingQuantity`.

```go
type Product struct {
    // ... existing fields ...
    ServingSize     string  `json:"serving_size"`      // e.g. "250ml", "1 bar (40g)"
    ServingQuantity float64 `json:"serving_quantity"`   // e.g. 250, 40
}
```

#### 6. ProductServicer Interface — Add `lang` parameter
**File**: [`product_controller.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/internal/controller/product_controller.go)
**Changes**: Update `ProductServicer` interface and handler methods.

```go
type ProductServicer interface {
    LookupBarcode(ctx context.Context, barcode string, lang string) (*types.Product, error)
    SearchProducts(ctx context.Context, query string, limit int, lang string) ([]types.Product, error)
}
```

Update `BarcodeHandler` and `SearchHandler` to extract `lang` from the input DTO and pass it through.

#### 7. Controller Input DTOs — Add `Lang` field
**File**: [`product_types.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/internal/controller/product_types.go)
**Changes**: Add `Lang` query parameter to both input types. Add `ServingSize` and `ServingQuantity` to `ProductBody`.

```go
type BarcodeInput struct {
    EAN  string `path:"ean" doc:"EAN/UPC barcode" example:"0049000000443"`
    Lang string `query:"lang,omitempty" doc:"Language code for localized results" example:"de"`
}

type SearchProductsInput struct {
    Query string `query:"query" required:"true" doc:"Search query" example:"yogurt"`
    Limit int    `query:"limit" default:"10" doc:"Max results to return" example:"10"`
    Lang  string `query:"lang,omitempty" doc:"Language code for localized results" example:"de"`
}

type ProductBody struct {
    // ... existing fields ...
    ServingSize     string  `json:"serving_size,omitempty"     doc:"Serving size description" example:"250ml"`
    ServingQuantity float64 `json:"serving_quantity,omitempty" doc:"Serving size quantity in grams" example:"250"`
}
```

#### 8. ProductService — Thread `lang` parameter
**File**: [`product_service.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/service/product_service.go)
**Changes**: Update `LookupBarcode` and `SearchProducts` to accept and pass `lang`.

```go
func (s *ProductService) LookupBarcode(ctx context.Context, barcode string, lang string) (*types.Product, error) {
    // ... pass lang to ds.LookupBarcode(ctx, barcode, lang) ...
}

func (s *ProductService) SearchProducts(ctx context.Context, query string, limit int, lang string) ([]types.Product, error) {
    // ... pass lang to ds.SearchProducts(ctx, query, remaining, lang) ...
}
```

#### 9. Log Endpoint — Make fields optional for product-based logging
**File**: [`nutrition_types.go` (controller)](file:///Users/dogab/code/MacroGuard/apps/api-go/internal/controller/nutrition_types.go)
**Changes**: Make `Confidence` a pointer and `Ingredients` use `omitempty`.

```go
type LogFoodInputBody struct {
    UserID      *string          `json:"user_id,omitempty" doc:"Optional UUID of the user"`
    FoodName    string           `json:"food_name" example:"Grilled Chicken Salad" doc:"Food name"`
    Confidence  *float64         `json:"confidence,omitempty" doc:"Detection confidence score (omit for product-based logs)"`
    Macros      *MacroData       `json:"macros" doc:"Nutritional macro information"`
    ServingSize string           `json:"serving_size,omitempty" doc:"Serving size description" example:"150g"`
    Ingredients []IngredientBody `json:"ingredients,omitempty" doc:"Breakdown of individual ingredients (optional for product logs)"`
}
```

**File**: [`nutrition_types.go` (service)](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/service/nutrition_types.go)
**Changes**: Make `Confidence` a pointer and `Ingredients` allow nil.

```go
type LogFoodInput struct {
    UserID      *string
    FoodName    string
    Confidence  *float64     // nil for product-based logs
    Macros      MacroData
    ServingSize string       // e.g. "150g"
    Ingredients []Ingredient // may be nil for product logs
}
```

**File**: [`nutrition_controller.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/internal/controller/nutrition_controller.go)
**Changes**: Update `LogFoodHandler` to handle nil `Confidence` (default to 1.0 for product-based logs) and empty `Ingredients`.

**File**: [`nutrition_service.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/service/nutrition_service.go)
**Changes**: Update `LogFood` to handle nil `Confidence` (default to 1.0) and nil `Ingredients` (create a single-ingredient log from the meal macros).

#### 10. Update All Tests
**Files**:
- [`product_service_test.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/service/product_service_test.go) — Update `mockDatasource` to accept `lang` parameter
- [`off_client_test.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/off_client_test.go) — Update call sites to pass `lang`
- [`fsvo_client_test.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/fsvo_client_test.go) — Update call sites to pass `lang`
- [`usda_client_test.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/usda_client_test.go) — Update call sites to pass `lang`
- Add new test: `TestOFFClient_LookupBarcode_PerRequestLang` — verify per-request `lang` overrides struct-level default
- Add new test: `TestFSVOClient_SearchProducts_PerRequestLang` — verify same for FSVO

#### 11. Regenerate OpenAPI Spec + TypeScript Client
```bash
make openapi   # runs gen-openapi-spec + gen-api-client
```

### Success Criteria:

#### Automated Verification:
- [ ] Code compiles: `cd apps/api-go && go build ./...`
- [ ] All tests pass: `cd apps/api-go && go test ./...`
- [ ] Formatting: `cd apps/api-go && gofmt -l .` reports no files
- [ ] TypeScript client regenerated: `make openapi` succeeds
- [ ] Frontend builds: `cd apps/web && pnpm check`

#### Manual Verification:
- [ ] `curl "http://localhost:8080/api/products/search?query=Yogurt&lang=de"` returns German product names from OFF
- [ ] `curl "http://localhost:8080/api/products/barcode/7613035466432?lang=de"` returns product with `serving_size` field populated
- [ ] `curl -X POST "http://localhost:8080/api/nutrition/log"` with product-shaped body (no confidence, no ingredients) succeeds

**⏸️ PAUSE: Manual confirmation required before proceeding to Phase 2.**

---

## Phase 2: Frontend — Barcode Library Evaluation

### Overview
Install and evaluate both candidate barcode scanning libraries to determine the best fit for the project.

### Changes Required:

#### 1. Install Both Libraries
```bash
cd apps/web && pnpm add html5-qrcode @niczvt/barcode-scanner-svelte
```

#### 2. Create Evaluation Test Harness
**File**: `apps/web/src/lib/components/food/BarcodeTestHarness.svelte` [NEW, TEMPORARY]
**Changes**: A simple Svelte 5 component that renders both scanners side-by-side with:
- Camera feed for each library
- EAN-13/UPC-A detection output
- Timing measurement (time from camera start to first successful decode)
- Bundle size comparison via build output

#### 3. Evaluation Criteria

| Criteria | Weight | html5-qrcode | @niczvt/barcode-scanner-svelte |
|----------|--------|--------------|------------------------------|
| EAN-13/UPC-A reliability | High | TBD | TBD |
| Svelte 5 compatibility | High | TBD | TBD |
| Bundle size | Medium | TBD | TBD |
| Camera handling quality | Medium | TBD | TBD |
| API ergonomics | Low | TBD | TBD |

### Success Criteria:

#### Manual Verification:
- [ ] Both libraries render camera feed successfully
- [ ] Test scanning EAN-13 barcode with physical product
- [ ] Compare bundle size impact: `pnpm build` with each library
- [ ] Document decision with rationale

**⏸️ PAUSE: Manual barcode testing required. Commit to chosen library, remove the other.**

---

## Phase 3: Frontend — Product Search & Barcode Scanner UI

### Overview
Build the full frontend flow: new modals, result cards, and the product logging drawer. Integrate into `AddEntryModal`.

### Changes Required:

#### 1. Update AddEntryModal — New Entry Points
**File**: [`AddEntryModal.svelte`](file:///Users/dogab/code/MacroGuard/apps/web/src/lib/components/navigation/AddEntryModal.svelte)
**Changes**: Add two new action cards below the existing three:
- **"Search Product"** (🔍 icon) — opens `ProductSearchModal`
- **"Scan Barcode"** (📷 icon) — opens `BarcodeScannerModal`

```svelte
<!-- New entries below existing Quick Add -->
<button onclick={() => { open = false; productSearchOpen = true; }}>
  <Search class="..." /> Search Product
</button>
<button onclick={() => { open = false; barcodeScannerOpen = true; }}>
  <ScanBarcode class="..." /> Scan Barcode
</button>
```

State management:
```svelte
let productSearchOpen = $state(false);
let barcodeScannerOpen = $state(false);
```

#### 2. ProductSearchModal [NEW]
**File**: `apps/web/src/lib/components/food/ProductSearchModal.svelte` [NEW]
**Purpose**: Full-screen modal with search input and product results.

Key implementation details:
- **Search input**: Text input with debounced API calls (300ms debounce)
- **API call**: `api.GET("/api/products/search", { params: { query: { query, limit: 15, lang: navigator.language.split("-")[0] } } })`
- **Results list**: Render `ProductResultCard` for each result
- **Loading state**: Skeleton cards while searching
- **Empty state**: "No products found" with suggestion to try different terms
- **Error state**: Toast/inline error for API failures
- **Selection**: Clicking a result card opens `ProductLogDrawer` with the selected product
- **Pattern**: Follow `FoodScannerModal.svelte` patterns for modal structure, transitions, and haptics

Props:
```typescript
interface Props {
  open: boolean;
}
```

Events:
- Dispatches `app:food-logged` custom event on successful log (same as FoodScannerModal)

#### 3. BarcodeScannerModal [NEW]
**File**: `apps/web/src/lib/components/food/BarcodeScannerModal.svelte` [NEW]
**Purpose**: Camera-based barcode scanner with automatic product lookup.

Key implementation details:
- **Camera**: Use evaluated barcode library (from Phase 2) to detect EAN-13/UPC-A
- **Auto-lookup**: On successful scan, immediately call `api.GET("/api/products/barcode/{ean}", { params: { path: { ean }, query: { lang: navigator.language.split("-")[0] } } })`
- **States**: `scanning` → `loading` (API lookup) → `found` (show product) / `not-found` (show error)
- **Found state**: Display product info via `ProductResultCard`, with "Log This Product" button opening `ProductLogDrawer`
- **Not found state**: "Product not found in our database" with a "Search Instead" button
- **Camera cleanup**: Stop camera stream on modal close (use `useCamera` hook pattern)

Props:
```typescript
interface Props {
  open: boolean;
}
```

#### 4. ProductResultCard [NEW]
**File**: `apps/web/src/lib/components/food/ProductResultCard.svelte` [NEW]
**Purpose**: Reusable card displaying a product's key info and macros.

Display:
- Product image (if available) or placeholder icon
- Product name + brand
- Source badge (OFF / FSVO / USDA) with subtle styling
- NutriScore badge (A-E colored circles, OFF only)
- Macro summary row: `56 kcal · 3.2g P · 5.1g C · 1.8g F` per 100g
- Serving size hint if available (e.g. "1 serving = 250ml")

Props:
```typescript
interface Props {
  product: ProductBody;
  onclick?: () => void;
}
```

#### 5. ProductLogDrawer [NEW]
**File**: `apps/web/src/lib/components/food/ProductLogDrawer.svelte` [NEW]
**Purpose**: Bottom drawer for selecting serving size and logging a product.

Key implementation details:
- **Serving options**:
  1. "100g" (default, pre-selected)
  2. Product serving size (if available from `serving_size` field, e.g. "1 bar (40g)")
  3. "Custom weight" — numeric input in grams
- **Macro computation**: Scale per-100g macros by `selectedWeight / 100`
- **Live preview**: Show computed macros updating as user changes serving size
- **Log button**: Call `api.POST("/api/nutrition/log", { body: { food_name, macros, serving_size } })` — no `confidence`, no `ingredients`
- **Success state**: Checkmark animation + auto-close
- **Dispatch**: `app:food-logged` event on success

Props:
```typescript
interface Props {
  open: boolean;
  product: ProductBody;
}
```

### shadcn-svelte Components Needed:

Check if these are installed; install any missing:
- `drawer` — for `ProductLogDrawer` (bottom sheet on mobile)
- `dialog` — for modals (likely already installed)
- `input` — for search input and custom weight
- `badge` — for source/NutriScore badges
- `skeleton` — for loading states

### Success Criteria:

#### Automated Verification:
- [ ] Frontend builds: `cd apps/web && pnpm check`
- [ ] No TypeScript errors
- [ ] No lint errors: `cd apps/web && pnpm lint`

#### Manual Verification:
- [ ] Open AddEntryModal → see 5 options including "Search Product" and "Scan Barcode"
- [ ] Search "Yogurt" → results load with macros displayed on cards
- [ ] Click a result → `ProductLogDrawer` opens with serving size options
- [ ] Select "Custom weight 250g" → macros update (2.5× the per-100g values)
- [ ] Click "Log" → success animation + daily intake counter updates
- [ ] Open BarcodeScannerModal → camera activates → scan a physical barcode → product loads
- [ ] Barcode not found → error state shows with "Search Instead" link

**⏸️ PAUSE: Full manual UI testing required before proceeding to Phase 4.**

---

## Phase 4: Architecture Documentation Updates

### Overview
Update all three architecture docs to reflect the new barcode scanning, product search, and per-request language features.

### Changes Required:

#### 1. Root Architecture
**File**: [`architecture.md`](file:///Users/dogab/code/MacroGuard/architecture.md)
**Changes**:
- Add barcode/search endpoints to the API Endpoints table:
  ```
  | GET  | /api/products/barcode/{ean} | Look up product by barcode |
  | GET  | /api/products/search        | Search products by name    |
  ```
- Update System Context diagram to mention barcode scanning

#### 2. Backend Architecture
**File**: [`apps/api-go/architecture.md`](file:///Users/dogab/code/MacroGuard/apps/api-go/architecture.md)
**Changes**:
- Document the `lang` query parameter on product endpoints
- Note the per-request language propagation pattern in the datasource section
- Add `serving_size`/`serving_quantity` to the Product type documentation
- Document the optional `confidence`/`ingredients` behavior on the log endpoint

#### 3. Frontend Architecture
**File**: [`apps/web/architecture.md`](file:///Users/dogab/code/MacroGuard/apps/web/architecture.md)
**Changes**:
- Add new components to the directory structure:
  ```
  │   │   ├── food/
  │   │   │   ├── FoodScannerModal.svelte    # AI food scan
  │   │   │   ├── ProductSearchModal.svelte  # Product search [NEW]
  │   │   │   ├── BarcodeScannerModal.svelte # Barcode scanner [NEW]
  │   │   │   ├── ProductResultCard.svelte   # Product card [NEW]
  │   │   │   ├── ProductLogDrawer.svelte    # Serving picker [NEW]
  │   │   │   └── ScanResultsDisplay.svelte  # AI scan results
  ```
- Document the product logging flow (search/scan → select → serving size → log)
- Note the chosen barcode library and its integration pattern

### Success Criteria:

#### Manual Verification:
- [ ] User confirms architecture documents accurately reflect the new implementation

**⏸️ PAUSE: Final manual confirmation before closing.**

---

## Testing Strategy

### Unit Tests (Backend):
- All existing tests updated with new `lang` parameter
- `TestOFFClient_LookupBarcode_PerRequestLang` — verify `lc=de` appears in OFF request URL
- `TestFSVOClient_SearchProducts_PerRequestLang` — verify `lang=fr` appears in FSVO request URL
- `TestProductService_LookupBarcode_CacheHit` — verify `lang` is passed through
- `TestLogFood_ProductBased` — verify log succeeds with nil `Confidence` and nil `Ingredients`

### Key Edge Cases:
- Empty `lang` parameter → falls back to struct-level default
- OFF product without `serving_size` → field is empty string, UI shows "100g" only
- Product log with 0g custom weight → validation error
- Barcode scan returns 404 → scanner shows "not found" + "search instead" link
- Search with special characters → URL-encoded safely

### Manual Testing Steps:
1. Start full stack: `make dev`
2. Open app → AddEntryModal → "Search Product" → search "Banana" → results load
3. Click a result → ProductLogDrawer opens → select "Custom weight 200g" → macros scale correctly
4. Click "Log" → success → daily intake reflects new entry
5. AddEntryModal → "Scan Barcode" → scan physical product → auto-lookup → log
6. Verify `/api/products/search?query=Yogurt&lang=de` returns German product names

## Performance Considerations

- **Debounced search**: 300ms debounce prevents excessive API calls during typing
- **Language param has no cache penalty**: Meilisearch cache is language-agnostic (stores the product as returned, subsequent lookups by barcode are cache hits regardless of language)
- **Barcode library bundle size**: Will be measured in Phase 2 evaluation — target < 100KB gzipped

## Migration Notes

None — no database schema changes. Product-based food logs use the same `food_logs` + `food_log_ingredients` tables. When `Ingredients` is nil, the log has just the top-level macros (no ingredient breakdown).

## References

- Research: [`thoughts/research/2026-04-05-8-barcode-scanning-search.md`](file:///Users/dogab/code/MacroGuard/thoughts/research/2026-04-05-8-barcode-scanning-search.md)
- GitHub Issue: #8
- Pattern reference: [`plan_fsvo-integration.md`](file:///Users/dogab/code/MacroGuard/thoughts/plan/plan_fsvo-integration.md)
- OFF API docs: https://wiki.openfoodfacts.org/API
