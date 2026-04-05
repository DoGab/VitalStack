---
date: "2026-04-04T20:49:00+02:00"
git_commit: 8585c669a7b3b0118f1100f4fe885d14624ef48b
branch: main
topic: "Integrate Swiss FSVO API as secondary datasource"
tags: [research, codebase, fsvo, datasource, product-search, backend]
status: complete
last_updated: "2026-04-04"
last_updated_note: "Resolved open question: FSVO language default confirmed as 'de'"
---

# Research: Integrate Swiss FSVO API as Secondary Datasource

**Date**: 2026-04-04T20:49:00+02:00
**Git Commit**: 8585c669a7b3b0118f1100f4fe885d14624ef48b
**Branch**: main

## Research Question

Integrate the Swiss FSVO (Federal Food Safety and Veterinary Office) API (`naehrwertdaten.ch`) as the **secondary** datasource in the product lookup waterfall, making USDA the **third** datasource. Document changes in architecture files. The FSVO client must use the same Meilisearch index for caching.

## Summary

The Swiss FSVO API is a free, authentication-free REST API with ~1,100 generic foods and full macro data (per 100g). It already maps cleanly onto the existing `FoodDatasource` interface. The integration is purely backend — a new `fsvo_client.go` in `pkg/datasource/`, config additions in `conf.go`, wiring in `server.go`, and a reorder of the datasource slice passed to `ProductService`.

## 🎯 Objective

Add the Swiss FSVO (naehrwertdaten.ch) REST API as a secondary datasource in the product search waterfall, between Open Food Facts (primary) and USDA (demoted to third). Products found via FSVO get cached in the shared Meilisearch index exactly like OFF and USDA products.

## 📋 Scope

- **Backend** — New FSVO client, config, wiring, reordered waterfall, tests
- **Architecture** — Update `apps/api-go/architecture.md` to document FSVO as secondary datasource

## ✅ Derived Definition of Done

- [ ] New `FSVOClient` in `pkg/datasource/fsvo_client.go` implements `FoodDatasource` interface
- [ ] `FSVOClient.SearchProducts()` calls FSVO `/foods` search, then enriches with `/food/{DBID}` for macros
- [ ] `FSVOClient.LookupBarcode()` returns `ErrNotFound` (FSVO has no barcodes)
- [ ] FSVO products are normalized to `types.Product` with `Source: "fsvo"` and `ID: "fsvo-{DBID}"`
- [ ] Waterfall order is: Meilisearch → OFF → **FSVO** → USDA
- [ ] FSVO config flags added to `conf.go` (base URL, language)
- [ ] FSVO client wired in `server.go` between OFF and USDA
- [ ] `local-config.yaml` updated with FSVO defaults
- [ ] `types/product.go` Source comment updated to include "fsvo"
- [ ] Architecture documentation updated in `apps/api-go/architecture.md`
- [ ] Unit tests for `fsvo_client.go` with httptest fixtures
- [ ] Existing tests still pass (`go test ./...`)

---

## Detailed Findings

### 1. Swiss FSVO API Structure

**Base URL:** `https://api.webapp.prod.blv.foodcase-services.com/BLV_WebApp_WS/webresources/BLV-api`
**Auth:** None (free, no API key required)
**Languages:** `de`, `en`, `fr`, `it` (query param `lang`)
**OpenAPI Spec:** [openapi.json](https://api.webapp.prod.blv.foodcase-services.com/BLV_WebApp_WS/webresources/openapi.json)

#### Key Endpoints

| Endpoint | Method | Purpose | Key Params |
|----------|--------|---------|------------|
| `/foods` | GET | Search foods by name | `search`, `lang`, `limit`, `offset`, `type` (generic) |
| `/food/{DBID}` | GET | Get single food with all nutrient values | `lang`, `DBID` (path) |
| `/components` | GET | List nutrient components (metadata) | `lang` |
| `/categories` | GET | List all categories | `lang` |

#### Search Response (`/foods?search=broccoli&lang=en&limit=2`)

```json
[
  {
    "id": 349687,
    "foodName": "Broccoli, raw",
    "generic": true,
    "categoryNames": "Fresh vegetables",
    "amount": 0.0,
    "foodid": 351,
    "valueTypeCode": ""
  },
  {
    "id": 350098,
    "foodName": "Broccoli, steamed (without addition of salt)",
    "generic": true,
    "categoryNames": "Cooked vegetables (incl. cans)",
    "amount": 0.0,
    "foodid": 1005,
    "valueTypeCode": ""
  }
]
```

**Key:** `id` is the `DBID` used for detail lookup. `foodName` is the localized name.

#### Food Detail Response (`/food/{DBID}`)

Returns full `Food` object with embedded `values[]` array containing all nutrients. Each value has:
- `value` (numeric, per 100g)
- `component.id` (maps to nutrient type)
- `component.code` (e.g., `"ENERCC"`, `"PROT625"`)
- `unit.code` (e.g., `"kcal"`, `"g"`)

#### FSVO Nutrient Component IDs for Macro Extraction

| Macro | Component ID | Code | Unit |
|-------|-------------|------|------|
| Energy (kcal) | `1777` | `ENERCC` | kcal |
| Protein | `15478` | `PROT625` | g |
| Carbs (available) | `55` | `CHO` | g |
| Fat (total) | `282` | `FAT` | g |
| Dietary fibres | `295` | `FIBT` | g |

**Verified:** Broccoli raw returns: 39 kcal, 3.6g protein, 3.2g carbs, 0.6g fat, 3.2g fiber.

### 2. Two-Step Lookup Pattern

The FSVO API requires a **two-step lookup** for search:
1. **Search** (`/foods?search=X`) → returns list of `{id, foodName, categoryNames}` (no macros)
2. **Detail** (`/food/{DBID}`) → returns full food with all nutrient `values[]`

For each search result, we need the detail call to get macros. This is unavoidable — the search endpoint only returns metadata.

**Implementation strategy:** For `SearchProducts()`, call `/foods` first, then batch the detail calls for the top N results. Keep it sequential since FSVO is a small dataset (~1,100 foods) and timeout concern is minimal.

### 3. No Barcode Support

FSVO has **no barcode data** — it's all generic foods (carrot, milk, bread). `LookupBarcode()` will immediately return `ErrNotFound`, which is the sentinel value that continues the waterfall.

### 4. Existing FoodDatasource Interface ([client.go:17-28](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/client.go#L17-L28))

```go
type FoodDatasource interface {
    LookupBarcode(ctx context.Context, barcode string) (*types.Product, error)
    SearchProducts(ctx context.Context, query string, limit int) ([]types.Product, error)
    Name() string
}
```

The FSVO client must implement all three methods. `LookupBarcode` returns `ErrNotFound` always. `SearchProducts` calls `/foods` then `/food/{DBID}` for each result. `Name()` returns `"fsvo"`.

### 5. Existing Waterfall in ProductService ([product_service.go:14-30](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/service/product_service.go#L14-L30))

```go
type ProductService struct {
    index       search.ProductSearchIndex
    datasources []datasource.FoodDatasource
}

func NewProductService(index search.ProductSearchIndex, datasources ...datasource.FoodDatasource) *ProductService {
```

The datasource order is determined by the variadic `datasources` argument order. Currently wired as `(meiliClient, offClient, usdaClient)` in [server.go:86](file:///Users/dogab/code/MacroGuard/apps/api-go/cmd/server.go#L86). To insert FSVO as secondary: `(meiliClient, offClient, fsvoClient, usdaClient)`.

### 6. Server Wiring ([server.go:69-89](file:///Users/dogab/code/MacroGuard/apps/api-go/cmd/server.go#L69-L89))

The FSVO client needs to be instantiated between OFF and USDA clients:

```go
// Current:
productSvc := service.NewProductService(meiliClient, offClient, usdaClient)

// After change:
fsvoClient := datasource.NewFSVOClient(http.DefaultClient, ...)
productSvc := service.NewProductService(meiliClient, offClient, fsvoClient, usdaClient)
```

### 7. Config System ([conf.go:128-149](file:///Users/dogab/code/MacroGuard/apps/api-go/internal/conf/conf.go#L128-L149))

New FSVO config constants needed:
- `fsvo.base-url` (default: `https://api.webapp.prod.blv.foodcase-services.com/BLV_WebApp_WS/webresources/BLV-api`)
- `fsvo.language` (default: `en`)

No API key needed — FSVO is free and unauthenticated.

### 8. Meilisearch Integration

FSVO products get cached in the **same Meilisearch index** (`products`) as OFF and USDA. The existing `cacheProductAsync` / `cacheProductsAsync` methods in `ProductService` handle this automatically. FSVO products will have IDs prefixed with `fsvo-` (e.g., `fsvo-349687`) to prevent collisions, following the existing pattern of `off-{barcode}` and `usda-{fdcId}`.

### 9. Test Patterns ([off_client_test.go](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/off_client_test.go), [usda_client_test.go](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/usda_client_test.go))

Both clients use:
- `httptest.NewServer` for mock HTTP servers
- JSON fixture files in `testdata/` directory
- `SetBaseURL()` method for test URL injection
- `loadTestFixture(t, "testdata/xxx.json")` helper (lives in `off_client_test.go`)
- Tests verify: success path, not-found, empty results, request path/params

The FSVO client tests should follow this exact pattern.

---

## Code References

- [client.go:17-28](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/client.go#L17-L28) — `FoodDatasource` interface
- [off_client.go:53-99](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/off_client.go#L53-L99) — OFFClient constructor with functional options pattern
- [usda_client.go:50-69](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/usda_client.go#L50-L69) — USDAClient constructor (simpler pattern)
- [product_service.go:23-30](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/service/product_service.go#L23-L30) — `NewProductService` with variadic datasources
- [server.go:69-89](file:///Users/dogab/code/MacroGuard/apps/api-go/cmd/server.go#L69-L89) — Product search wiring (where FSVO must be inserted)
- [conf.go:119-149](file:///Users/dogab/code/MacroGuard/apps/api-go/internal/conf/conf.go#L119-L149) — USDA and OFF config patterns to follow
- [product_service_test.go:1-225](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/service/product_service_test.go) — Service-level waterfall tests
- [usda_client_test.go:1-164](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/usda_client_test.go) — Client test patterns
- [architecture.md:1-204](file:///Users/dogab/code/MacroGuard/apps/api-go/architecture.md) — Backend architecture to update

## 🔗 Architecture Documentation

- Backend: `apps/api-go/architecture.md` (needs update for FSVO)
- Frontend: `apps/web/architecture.md` (no changes needed — backend-only feature)

## Files to Create

| File | Purpose |
|------|---------|
| `pkg/datasource/fsvo_client.go` | FSVO HTTP client implementing `FoodDatasource` |
| `pkg/datasource/fsvo_client_test.go` | Unit tests with httptest |
| `pkg/datasource/testdata/fsvo_search.json` | Search fixture |
| `pkg/datasource/testdata/fsvo_food.json` | Food detail fixture |

## Files to Modify

| File | Change |
|------|--------|
| `internal/conf/conf.go` | Add FSVO config constants (`fsvo.base-url`, `fsvo.language`) |
| `cmd/server.go` | Wire FSVOClient between OFF and USDA |
| `local-config.yaml` | Add FSVO config defaults |
| `pkg/types/product.go` | Update Source comment to include `"fsvo"` |
| `apps/api-go/architecture.md` | Document FSVO as secondary datasource |

## Waterfall Order (After Change)

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
  │  2. OPEN      │  ← Primary: widest branded product coverage
  │  FOOD FACTS   │     Barcode + text search
  └──────┬───────┘
         │ miss or generic food
         ▼
  ┌──────────────┐
  │  3. SWISS     │  ← Secondary: Swiss reference generic foods
  │    FSVO       │     ~1,100 gov-verified foods, multilingual
  └──────┬───────┘
         │ miss
         ▼
  ┌──────────────┐
  │  4. USDA      │  ← Tertiary: US raw/generic foods
  │     FDC       │     450K+ items, lab-analyzed
  └──────────────┘
```

## Historical Context (from thoughts/)

- `thoughts/research/research_food-datasources.md` — Original datasource comparison that identified FSVO as a "Phase 2 enrichment source" with ~1,100 generic foods

## ✅ Resolved Questions

- **Language default:** ✅ Confirmed — default FSVO language is `"de"` (German, Swiss-primary). Configurable via `fsvo.language` config flag.

---

> **Research complete. Ready for Planning phase.**
