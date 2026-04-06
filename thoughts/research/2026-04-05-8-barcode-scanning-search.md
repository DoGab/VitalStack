---
issue: "#8"
title: "Support Barcode scanning and product search"
date: 2026-04-05
tags: [research, codebase, barcode, product-search, frontend, backend, i18n]
status: complete
last_updated: 2026-04-06
last_updated_note: "Resolved open questions with user decisions"
---

# Research: Issue #8 — Support Barcode Scanning and Product Search

> **GitHub Issue:** [#8 — Support Barcode scanning and product search](https://github.com/DoGab/VitalStack/issues/8)
> **Date:** 2026-04-05

## Research Question

How do we extend VitalStack to support barcode scanning and product-text-search from the frontend, while propagating the user's localized language through the backend waterfall architecture to the external datasource clients?

## Executive Summary

The backend is **fully functional** — product search and barcode lookup APIs already exist at `/api/products/search` and `/api/products/barcode/{ean}`. The four-tier waterfall (Meilisearch → OFF → FSVO → USDA) works end-to-end. The **missing pieces** are:

1. **Backend:** Language is hardcoded at startup via `viper` config. The API must accept a `lang` query parameter and propagate it per-request to OFF/FSVO clients.
2. **Frontend:** The `AddEntryModal` only offers "Take Photo", "Upload Image", and "Quick Add". It needs two new entry points: **"Search Product"** and **"Scan Barcode"**. A new results UI is needed to display products with macro data and logging capabilities.

---

## 1. Backend Architecture

### 1.1. Current Waterfall Strategy

[product_service.go:14-30](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/service/product_service.go#L14-L30) — `ProductService` struct and constructor.

```
User Request (Barcode or Search Query)
           │
  ┌────────▼───────┐
  │  1. MEILISEARCH │  ← Local cache (products index)
  │     (Cache)     │     If enough results → return immediately
  └───────┬────────┘
          │ miss
  ┌───────▼────────┐
  │  2. OPEN FOOD  │  ← Primary: widest branded product coverage
  │     FACTS      │     Barcode + text search, retries on 50x
  └───────┬────────┘
          │ miss
  ┌───────▼────────┐
  │  3. SWISS FSVO │  ← Generic foods (~1,100), no barcodes
  │  (naehrwert)   │     Two-step: search → fetch detail for macros
  └───────┬────────┘
          │ miss
  ┌───────▼────────┐
  │  4. USDA       │  ← Tertiary: US generic/branded foods
  │  FoodData      │     No barcode support either
  └────────────────┘
```

### 1.2. FoodDatasource Interface

[client.go:17-28](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/client.go#L17-L28)

```go
type FoodDatasource interface {
    LookupBarcode(ctx context.Context, barcode string) (*types.Product, error)
    SearchProducts(ctx context.Context, query string, limit int) ([]types.Product, error)
    Name() string
}
```

**Key observation:** The interface does **not** accept a `language` parameter. Language is set at client construction time via functional options (`WithLanguage`, `WithFSVOLanguage`).

### 1.3. Language Configuration (Current State)

| Client     | Option             | Default | Config Key              | Set In                                                                                              |
|------------|--------------------|---------|-------------------------|-----------------------------------------------------------------------------------------------------|
| OFFClient  | `WithLanguage()`   | `"en"`  | `conf.OFFLanguageArg`   | [server.go:82](file:///Users/dogab/code/MacroGuard/apps/api-go/cmd/server.go#L82)                   |
| FSVOClient | `WithFSVOLanguage()` | `"de"` | `conf.FSVOLanguageArg`  | [server.go:88](file:///Users/dogab/code/MacroGuard/apps/api-go/cmd/server.go#L88)                   |
| USDAClient | N/A                | N/A     | N/A                     | USDA has no language support                                                                         |

**The Problem:** Both `OFFClient.language` and `FSVOClient.language` are struct-level fields set once at bootup. To support per-request language, we have two options:

1. **Option A — Interface change:** Add `language` to `FoodDatasource.SearchProducts()` and `LookupBarcode()` signatures.
2. **Option B — Context-based:** Store language in `context.Context` and read it in the clients.

**Recommendation: Option A (Interface change)** — It's explicit, testable, and the interface has only 3 implementors, so the blast radius is small.

### 1.4. Product Controller & Types

[product_controller.go:56-97](file:///Users/dogab/code/MacroGuard/apps/api-go/internal/controller/product_controller.go#L56-L97) — Handlers for barcode and search.

[product_types.go:10-14](file:///Users/dogab/code/MacroGuard/apps/api-go/internal/controller/product_types.go#L10-L14) — `SearchProductsInput`:

```go
type SearchProductsInput struct {
    Query string `query:"query" required:"true" doc:"Search query" example:"yogurt"`
    Limit int    `query:"limit" default:"10" doc:"Max results to return" example:"10"`
}
```

**Action Required:** Add `Lang string \`query:"lang" ...\`` to `SearchProductsInput` and `BarcodeInput`.

### 1.5. ProductServicer Interface

[product_controller.go:19-22](file:///Users/dogab/code/MacroGuard/apps/api-go/internal/controller/product_controller.go#L19-L22):

```go
type ProductServicer interface {
    LookupBarcode(ctx context.Context, barcode string) (*types.Product, error)
    SearchProducts(ctx context.Context, query string, limit int) ([]types.Product, error)
}
```

**Action Required:** Add `lang string` parameter to both methods.

---

## 2. Frontend Architecture

### 2.1. Component Tree (Entry Flow)

```
MobileDock.svelte  ──(+ button)──►  AddEntryModal.svelte
                                       │
                ┌──────────────────────┼──────────────────────┐
                │                      │                      │
        "Take Photo"           "Upload Image"          "Quick Add"
                │                      │
                └──────┬───────────────┘
                       ▼
             FoodScannerModal.svelte
                       │
                       ▼
            ScanResultsDisplay.svelte
```

### 2.2. AddEntryModal — Entry Point

[AddEntryModal.svelte](file:///Users/dogab/code/MacroGuard/apps/web/src/lib/components/navigation/AddEntryModal.svelte)

- Uses `Dialog` (desktop) / `Drawer` (mobile) pattern via `IsMobile` hook
- Options list is defined as a const array `addOptions` — **easy to extend**
- `handleAddOption()` dispatches to `FoodScannerModal` — needs new cases for `"barcode"` and `"search"`

### 2.3. FoodScannerModal — Camera + AI Scan

[FoodScannerModal.svelte](file:///Users/dogab/code/MacroGuard/apps/web/src/lib/components/food/FoodScannerModal.svelte) — 491 lines

- Handles camera capture, file upload, AI food scanning, and logging
- Already uses `useCamera` hook, which could be reused for barcode scanning
- The component is complex and self-contained — **barcode scanning should be a separate component**

### 2.4. ScanResultsDisplay — Macro Display

[ScanResultsDisplay.svelte](file:///Users/dogab/code/MacroGuard/apps/web/src/lib/components/food/ScanResultsDisplay.svelte) — 196 lines

- Displays macro breakdown with `CircularProgress` + `MacroBars` components
- Supports `preview` and `details` modes
- **Can be partially reused for product search results** (the macro display card pattern)

### 2.5. API Client (Frontend)

[client.ts](file:///Users/dogab/code/MacroGuard/apps/web/src/lib/api/client.ts) — Type-safe `openapi-fetch` client

[schema.d.ts](file:///Users/dogab/code/MacroGuard/apps/web/src/lib/api/schema.d.ts) — Auto-generated from `openapi.yaml`

Product-related operations already typed:
- `operations["lookup-product-barcode"]` → `GET /api/products/barcode/{ean}`
- `operations["search-products"]` → `GET /api/products/search?query=...&limit=...`

**Schema types available:**
- `ProductBody` — id, barcode, name, brand, image_url, source, nutri_score, macros
- `MacrosPer100gBody` — calories, protein, carbs, fat, fiber
- `SearchProductsOutputBody` — products[], attribution

### 2.6. Installed UI Components (shadcn-svelte)

Available in `apps/web/src/lib/components/ui/`:

| Component      | Usage               |
|----------------|---------------------|
| `dialog`       | Desktop modals      |
| `drawer`       | Mobile modals       |
| `card`         | Content containers  |
| `button`       | Actions             |
| `input`        | Text fields         |
| `badge`        | Tags/indicators     |
| `accordion`    | Collapsible sections|
| `skeleton`     | Loading states      |
| `progress`     | Linear progress     |
| `sheet`        | Slide-over panels   |
| `separator`    | Dividers            |

### 2.7. Hooks & State

| File                     | Purpose                          |
|--------------------------|----------------------------------|
| `is-mobile.svelte.ts`   | Responsive breakpoint detection  |
| `use-camera.svelte.ts`  | Camera stream lifecycle          |
| `nutrition.svelte.ts`   | Global nutrition state (runes)   |

---

## 3. Visual Baseline

### 3.1. Dashboard (Current)

![Dashboard baseline](file:///Users/dogab/.gemini/antigravity/brain/71bd9704-9e0f-490f-b3df-19c015219036/main_dashboard_page_1775409220552.png)

### 3.2. Add Entry Modal (Current)

![Add Entry modal baseline](file:///Users/dogab/.gemini/antigravity/brain/71bd9704-9e0f-490f-b3df-19c015219036/add_entry_modal_1775409225220.png)

The modal currently shows 3 options. Issue #8 requires adding **"Search Product"** and **"Scan Barcode"** (possibly as a tab within the camera view).

---

## 4. Implementation Scope

### 4.1. Backend Changes

| File | Change |
|------|--------|
| `pkg/datasource/client.go` | Add `lang string` param to `FoodDatasource.SearchProducts()` and `LookupBarcode()` |
| `pkg/datasource/off_client.go` | Accept `lang` param in methods, override struct-level language per-request |
| `pkg/datasource/fsvo_client.go` | Accept `lang` param in methods, override struct-level language per-request |
| `pkg/datasource/usda_client.go` | Accept `lang` param (ignored, USDA has no i18n) |
| `pkg/service/product_service.go` | Thread `lang` from service methods to datasource calls |
| `internal/controller/product_types.go` | Add `Lang` field to `SearchProductsInput` and `BarcodeInput` |
| `internal/controller/product_controller.go` | Pass `input.Lang` to service methods |
| `docs/openapi.yaml` | Add `lang` query parameter to both product endpoints |
| All `*_test.go` files | Update test signatures |

### 4.2. Frontend Changes

| File | Change |
|------|--------|
| `navigation/AddEntryModal.svelte` | Add "Search Product" and "Scan Barcode" options |
| `food/ProductSearchModal.svelte` | **[NEW]** — Search input + debounced API call + results list |
| `food/BarcodeScannerModal.svelte` | **[NEW]** — Camera barcode detection + auto-lookup |
| `food/ProductResultCard.svelte` | **[NEW]** — Product card with image, name, macros, + button |
| `food/ProductLogDrawer.svelte` | **[NEW]** — Serving size picker + log confirmation |
| `api/schema.d.ts` | Regenerated via `make openapi` after backend changes |

### 4.3. New Dependencies (Frontend)

- **Barcode scanning library:** `html5-qrcode` or `@niczvt/barcode-scanner-svelte` — needs evaluation. `html5-qrcode` is mature and framework-agnostic with strong EAN-13/UPC-A support.

---

## 5. Key References

### Backend
- [product_service.go](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/service/product_service.go) — Waterfall orchestration
- [product_controller.go](file:///Users/dogab/code/MacroGuard/apps/api-go/internal/controller/product_controller.go) — HTTP handlers
- [product_types.go](file:///Users/dogab/code/MacroGuard/apps/api-go/internal/controller/product_types.go) — Request/response DTOs
- [client.go](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/client.go) — FoodDatasource interface
- [off_client.go](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/off_client.go) — OFF client with language config
- [fsvo_client.go](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/fsvo_client.go) — FSVO client with language config
- [server.go](file:///Users/dogab/code/MacroGuard/apps/api-go/cmd/server.go) — Dependency wiring

### Frontend
- [AddEntryModal.svelte](file:///Users/dogab/code/MacroGuard/apps/web/src/lib/components/navigation/AddEntryModal.svelte) — Entry point for food logging
- [FoodScannerModal.svelte](file:///Users/dogab/code/MacroGuard/apps/web/src/lib/components/food/FoodScannerModal.svelte) — Camera + AI scan flow
- [ScanResultsDisplay.svelte](file:///Users/dogab/code/MacroGuard/apps/web/src/lib/components/food/ScanResultsDisplay.svelte) — Macro display
- [MobileDock.svelte](file:///Users/dogab/code/MacroGuard/apps/web/src/lib/components/navigation/MobileDock.svelte) — Bottom nav with + button
- [schema.d.ts](file:///Users/dogab/code/MacroGuard/apps/web/src/lib/api/schema.d.ts) — Generated API types

### Prior Research
- [research_food-datasources.md](file:///Users/dogab/code/MacroGuard/thoughts/research/research_food-datasources.md) — Original datasource comparison
- [research_fsvo-integration.md](file:///Users/dogab/code/MacroGuard/thoughts/research/research_fsvo-integration.md) — FSVO integration details

---

## 6. Open Questions

## 6. Resolved Decisions

1. **Barcode Scanner Library:** `html5-qrcode` vs `@niczvt/barcode-scanner-svelte` — **TBD, requires evaluation** during implementation planning. Both are candidates; evaluate bundle size, EAN-13/UPC-A reliability, and Svelte 5 compatibility before choosing.
2. **Barcode UI placement:** ✅ **Standalone `BarcodeScannerModal`** — separate from `FoodScannerModal` for clean separation of concerns. Not a tab within the camera view.
3. **Product logging flow:** ✅ **Use a serving-size picker.** When a user taps "+" on a product search result, present a `ProductLogDrawer` with serving-size options (e.g., "100g", "1 serving", custom weight) before logging. Do NOT log per-100g directly.
4. **Language detection:** ✅ **Use `navigator.language`** and extract the base language code (e.g., `"de-CH"` → `"de"`). User-configurable language preference in Profile deferred to a later stage.

---

## 7. Follow-up Research (2026-04-06)

### Barcode Library Evaluation (deferred to planning)

The user did not express a preference between `html5-qrcode` and `@niczvt/barcode-scanner-svelte`. During the planning phase, evaluate:

| Criterion | `html5-qrcode` | `@niczvt/barcode-scanner-svelte` |
|-----------|-----------------|----------------------------------|
| Maturity  | Widely used, framework-agnostic | Svelte-native wrapper |
| EAN-13/UPC-A support | Native | Depends on underlying engine |
| Bundle size | ~150KB gzipped | Needs measurement |
| Svelte 5 compat | Manual integration | May need verification |
| Camera API | Own implementation | Wraps browser API |

**Action for planning phase:** Install both locally, test with a physical EAN-13 barcode, and measure bundle impact before committing.

---

> **Research complete. Ready for Planning phase.**

