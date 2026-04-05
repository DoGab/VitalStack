# Swiss FSVO Datasource Integration — Implementation Plan

## Overview

Add the Swiss Federal Food Safety and Veterinary Office (FSVO / naehrwertdaten.ch) REST API as the **secondary** datasource in the product search waterfall. USDA moves to **third** position. FSVO provides ~1,100 government-verified generic food compositions with full macro data per 100g. The API is free, unauthenticated, and multilingual.

## Current State Analysis

The product search uses a waterfall pattern (`Meilisearch → OFF → USDA`) orchestrated by `ProductService`. The `FoodDatasource` interface is well-established and accepts variadic datasources in order of priority.

### Key Discoveries:
- [`FoodDatasource` interface](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/client.go#L17-L28) — three methods: `LookupBarcode`, `SearchProducts`, `Name`
- [`USDAClient`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/usda_client.go) — cleanest pattern to follow (simple constructor, `SetBaseURL()` for tests, nutrient ID switch)
- [Server wiring](file:///Users/dogab/code/MacroGuard/apps/api-go/cmd/server.go#L86) — datasource order set by argument position
- [Config system](file:///Users/dogab/code/MacroGuard/apps/api-go/internal/conf/conf.go#L119-L149) — Viper constants per concern with `Arg/Default/Help` triples
- [Product type](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/types/product.go#L13) — `Source` field comment needs `"fsvo"` added
- FSVO API requires **two-step lookup**: `/foods?search=X` (metadata only) → `/food/{DBID}` (full nutrients)
- FSVO has **no barcode data** — `LookupBarcode()` always returns `ErrNotFound`

## Desired End State

After implementation:
1. Searching "Broccoli" in the app hits FSVO and returns Swiss-verified macro data
2. FSVO results are cached in Meilisearch with `fsvo-{DBID}` IDs
3. Waterfall order: **Meilisearch → OFF → FSVO → USDA**
4. `architecture.md` documents the four-tier waterfall
5. All existing tests pass + new FSVO-specific tests pass

## What We're NOT Doing

- No frontend changes (backend-only feature)
- No FSVO barcode lookup (the API has no barcode data)
- No batch preloading of FSVO data into Meilisearch (on-demand caching only)
- No FSVO category browsing (search only, same as OFF/USDA)
- No new API endpoints — existing `/api/products/search` and `/api/products/barcode/{ean}` automatically pick up FSVO via waterfall

## Implementation Approach

Follow the patterns established by `OFFClient` (functional options) and `USDAClient` (nutrient switch):
1. Define response structs matching FSVO JSON
2. Functional option pattern: `NewFSVOClient(httpClient, baseURL, opts...)` with `WithFSVOLanguage()`
3. `SetBaseURL()` for test injection
4. Nutrient extraction via component code switch (`ENERCC`, `PROT625`, etc.)
5. Two-step search: call `/foods`, then `/food/{DBID}` per result for macros
6. No hardcoded default URL — defaults are managed exclusively by the `conf` package

---

## Phase 1: FSVO Client + Config + Wiring

### Overview
Implement the full FSVO client, add configuration, and wire it into the waterfall.

### Changes Required:

#### 1. Configuration Constants
**File**: `apps/api-go/internal/conf/conf.go`
**Changes**: Add FSVO config block between USDA and OFF sections.

```go
// Swiss FSVO (naehrwertdaten.ch)
fsvoKey = "fsvo."
// FSVOBaseURLArg is the flag name for the FSVO base URL
FSVOBaseURLArg = fsvoKey + "base-url"
// FSVOBaseURLDefault is the default value for the FSVO base URL
FSVOBaseURLDefault = "https://api.webapp.prod.blv.foodcase-services.com/BLV_WebApp_WS/webresources/BLV-api"
// FSVOBaseURLHelp is the help message for the FSVO base URL
FSVOBaseURLHelp = "Swiss FSVO food composition database base URL"

// FSVOLanguageArg is the flag name for the FSVO language
FSVOLanguageArg = fsvoKey + "language"
// FSVOLanguageDefault is the default value for the FSVO language
FSVOLanguageDefault = "de"
// FSVOLanguageHelp is the help message for the FSVO language
FSVOLanguageHelp = "FSVO response language code (de, en, fr, it)"
```

In `RegisterFlags`, add after the USDA block:
```go
// FSVO
pflags.String(FSVOBaseURLArg, FSVOBaseURLDefault, FSVOBaseURLHelp)
pflags.String(FSVOLanguageArg, FSVOLanguageDefault, FSVOLanguageHelp)
```

#### 2. FSVO Client Implementation
**File**: `apps/api-go/pkg/datasource/fsvo_client.go` [NEW]
**Changes**: New file implementing `FoodDatasource`.

Key implementation details:
- **Response structs**: `fsvoSearchFood` (from `/foods`) and `fsvoFood` / `fsvoValue` / `fsvoComponent` (from `/food/{DBID}`)
- **Functional options**: `FSVOClientOption` type with `WithFSVOLanguage()` — mirrors `OFFClientOption` pattern
- **Constructor**: `NewFSVOClient(httpClient *http.Client, baseURL string, opts ...FSVOClientOption) *FSVOClient`
- **`LookupBarcode`**: returns `ErrNotFound` immediately
- **`SearchProducts`**: calls `/foods?search=X&lang=Y&limit=Z`, then `/food/{DBID}` for each result
- **`Name()`**: returns `"fsvo"`
- **Nutrient extraction**: switch on `component.Code` matching `ENERCC`, `PROT625`, `CHO`, `FAT`, `FIBT`
- **No hardcoded default URL**: the `baseURL` is always passed from `conf.FSVOBaseURLDefault` via viper

```go
package datasource

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strconv"

    "github.com/dogab/vitalstack/api/pkg/types"
)

const (
    fsvoSource = "fsvo"

    // FSVO component codes for the macros we track.
    fsvoCodeEnergy  = "ENERCC"
    fsvoCodeProtein = "PROT625"
    fsvoCodeCarbs   = "CHO"
    fsvoCodeFat     = "FAT"
    fsvoCodeFiber   = "FIBT"
)

// --- FSVO API response types ---

// fsvoSearchFood represents a single item from the /foods search endpoint.
type fsvoSearchFood struct {
    ID            int    `json:"id"`       // DBID used for detail lookup
    FoodName      string `json:"foodName"`
    CategoryNames string `json:"categoryNames"`
}

// fsvoFood represents the full food detail from /food/{DBID}.
type fsvoFood struct {
    Name       string         `json:"name"`
    ID         int            `json:"id"`
    Categories []fsvoCategory `json:"categories"`
    Values     []fsvoValue    `json:"values"`
}

// fsvoCategory represents a food category.
type fsvoCategory struct {
    Name string `json:"name"`
    ID   int    `json:"id"`
}

// fsvoValue represents a single nutrient value entry.
type fsvoValue struct {
    Value     float64        `json:"value"`
    Component fsvoComponent  `json:"component"`
}

// fsvoComponent identifies what nutrient a value represents.
type fsvoComponent struct {
    Name string `json:"name"`
    ID   int    `json:"id"`
    Code string `json:"code"`
}

// --- Client ---

// FSVOClient is an HTTP client for the Swiss FSVO food composition database.
type FSVOClient struct {
    httpClient *http.Client
    baseURL    string
    language   string
}

// FSVOClientOption defines a functional option for configuring the FSVOClient.
type FSVOClientOption func(*FSVOClient)

// WithFSVOLanguage configures the language code (e.g. "de") for FSVO API requests.
// This allows overriding the default language per-client, enabling future
// per-request language support from the user's client.
func WithFSVOLanguage(lang string) FSVOClientOption {
    return func(c *FSVOClient) {
        c.language = lang
    }
}

// NewFSVOClient creates a new FSVO client with functional options.
// baseURL is required and should be provided from the conf package.
func NewFSVOClient(httpClient *http.Client, baseURL string, opts ...FSVOClientOption) *FSVOClient {
    c := &FSVOClient{
        httpClient: httpClient,
        baseURL:    baseURL,
        language:   "de",
    }

    for _, opt := range opts {
        opt(c)
    }

    return c
}

// SetBaseURL overrides the base URL (used for testing with httptest).
func (c *FSVOClient) SetBaseURL(u string) {
    c.baseURL = u
}

// Name returns the datasource identifier.
func (c *FSVOClient) Name() string {
    return fsvoSource
}

// LookupBarcode always returns ErrNotFound because FSVO has no barcode data.
func (c *FSVOClient) LookupBarcode(_ context.Context, _ string) (*types.Product, error) {
    return nil, ErrNotFound
}

// SearchProducts searches the FSVO database by name.
// This is a two-step process: search for foods, then fetch detail for each.
func (c *FSVOClient) SearchProducts(ctx context.Context, query string, limit int) ([]types.Product, error) {
    // Step 1: Search for matching foods (returns metadata only, no macros).
    searchFoods, err := c.searchFoods(ctx, query, limit)
    if err != nil {
        return nil, err
    }

    // Step 2: Fetch full detail for each search result to get macros.
    products := make([]types.Product, 0, len(searchFoods))
    for _, sf := range searchFoods {
        food, err := c.getFood(ctx, sf.ID)
        if err != nil {
            // Skip individual failures — don't fail the whole batch.
            continue
        }
        products = append(products, fsvoFoodToDomain(*food))
    }

    return products, nil
}

// searchFoods calls the FSVO /foods search endpoint.
func (c *FSVOClient) searchFoods(ctx context.Context, query string, limit int) ([]fsvoSearchFood, error) {
    params := url.Values{
        "search": {query},
        "lang":   {c.language},
        "limit":  {strconv.Itoa(limit)},
    }

    reqURL := fmt.Sprintf("%s/foods?%s", c.baseURL, params.Encode())

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
    if err != nil {
        return nil, fmt.Errorf("fsvo: creating search request: %w", err)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("fsvo: executing search request: %w", err)
    }
    defer func() { _ = resp.Body.Close() }()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("fsvo: search unexpected status %d", resp.StatusCode)
    }

    var foods []fsvoSearchFood
    if err := json.NewDecoder(resp.Body).Decode(&foods); err != nil {
        return nil, fmt.Errorf("fsvo: decoding search response: %w", err)
    }

    return foods, nil
}

// getFood fetches a single food by DBID from the FSVO /food/{DBID} endpoint.
func (c *FSVOClient) getFood(ctx context.Context, dbid int) (*fsvoFood, error) {
    reqURL := fmt.Sprintf("%s/food/%d?lang=%s", c.baseURL, dbid, c.language)

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
    if err != nil {
        return nil, fmt.Errorf("fsvo: creating food request: %w", err)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("fsvo: executing food request: %w", err)
    }
    defer func() { _ = resp.Body.Close() }()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("fsvo: food unexpected status %d", resp.StatusCode)
    }

    var food fsvoFood
    if err := json.NewDecoder(resp.Body).Decode(&food); err != nil {
        return nil, fmt.Errorf("fsvo: decoding food response: %w", err)
    }

    return &food, nil
}

// fsvoFoodToDomain converts an FSVO food detail to the domain Product type.
func fsvoFoodToDomain(f fsvoFood) types.Product {
    var categories []string
    for _, cat := range f.Categories {
        if cat.Name != "" {
            categories = append(categories, cat.Name)
        }
    }

    return types.Product{
        ID:         fmt.Sprintf("fsvo-%d", f.ID),
        Name:       f.Name,
        Categories: categories,
        Source:     fsvoSource,
        Macros:     extractFSVOMacros(f.Values),
    }
}

// extractFSVOMacros extracts macronutrient values from FSVO nutrient data by matching component codes.
func extractFSVOMacros(values []fsvoValue) types.MacrosPer100g {
    var macros types.MacrosPer100g
    for _, v := range values {
        switch v.Component.Code {
        case fsvoCodeEnergy:
            macros.Calories = v.Value
        case fsvoCodeProtein:
            macros.Protein = v.Value
        case fsvoCodeCarbs:
            macros.Carbs = v.Value
        case fsvoCodeFat:
            macros.Fat = v.Value
        case fsvoCodeFiber:
            macros.Fiber = v.Value
        }
    }
    return macros
}
```

#### 3. Server Wiring
**File**: `apps/api-go/cmd/server.go`
**Changes**: Instantiate `FSVOClient` between OFF and USDA, pass to `ProductService`.

```go
// Line ~84: after offClient, before usdaClient
fsvoClient := datasource.NewFSVOClient(
    http.DefaultClient,
    viper.GetString(conf.FSVOBaseURLArg),
    datasource.WithFSVOLanguage(viper.GetString(conf.FSVOLanguageArg)),
)
usdaClient := datasource.NewUSDAClient(http.DefaultClient, viper.GetString(conf.USDAAPIKeyArg))
productSvc := service.NewProductService(meiliClient, offClient, fsvoClient, usdaClient)
```

#### 4. Local Config
**File**: `apps/api-go/local-config.yaml`
**Changes**: Add FSVO section between `usda:` and `openfoodfacts:`.

```yaml
fsvo:
  base-url: "https://api.webapp.prod.blv.foodcase-services.com/BLV_WebApp_WS/webresources/BLV-api"
  language: "de"
```

#### 5. Product Type Source Comment
**File**: `apps/api-go/pkg/types/product.go`
**Changes**: Update the `Source` field doc comment to include `"fsvo"`.

```go
Source     string        `json:"source"`      // "openfoodfacts", "fsvo", "usda"
```

### Success Criteria:

#### Automated Verification:
- [x] Code compiles: `cd apps/api-go && go build ./...`
- [x] Existing tests still pass: `cd apps/api-go && go test ./...`
- [x] Formatting is correct: `cd apps/api-go && gofmt -l .` reports no files

#### Manual Verification:
- [x] Start the server with `make dev-api`, search "Broccoli" → FSVO results appear with correct macros (39 kcal, 3.6g protein)
- [x] Barcode lookup still works (FSVO is skipped, OFF/USDA respond)

**Implementation Note**: After completing this phase and all automated verification passes, pause here for manual confirmation.

---

## Phase 2: Unit Tests

### Overview
Add comprehensive unit tests for the FSVO client following the established httptest patterns.

### Changes Required:

#### 1. Search Fixture
**File**: `apps/api-go/pkg/datasource/testdata/fsvo_search.json` [NEW]
**Changes**: Fixture for the `/foods` search endpoint.

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

#### 2. Food Detail Fixture
**File**: `apps/api-go/pkg/datasource/testdata/fsvo_food.json` [NEW]
**Changes**: Fixture for the `/food/{DBID}` endpoint (subset of Broccoli raw).

```json
{
  "name": "Broccoli, raw",
  "id": 349687,
  "isgeneric": true,
  "isrecipe": false,
  "synonyms": [],
  "categories": [
    { "name": "Fresh vegetables", "id": 6637 }
  ],
  "values": [
    {
      "id": 6268031,
      "value": 39,
      "component": { "name": "Energy, kilocalories", "id": 1777, "code": "ENERCC" }
    },
    {
      "id": 6268030,
      "value": 3.6,
      "component": { "name": "Protein", "id": 15478, "code": "PROT625" }
    },
    {
      "id": 6268032,
      "value": 3.2,
      "component": { "name": "Carbohydrates, available", "id": 55, "code": "CHO" }
    },
    {
      "id": 6268061,
      "value": 0.6,
      "component": { "name": "Fat, total", "id": 282, "code": "FAT" }
    },
    {
      "id": 6268053,
      "value": 3.2,
      "component": { "name": "Dietary fibres", "id": 295, "code": "FIBT" }
    }
  ]
}
```

#### 3. Test File
**File**: `apps/api-go/pkg/datasource/fsvo_client_test.go` [NEW]
**Changes**: Unit tests covering search success, search empty, barcode not-found, macro mapping, and request path/params.

The test server needs to handle **two endpoints**: `/foods` (search) and `/food/{DBID}` (detail). Use URL path routing in the handler.

```go
package datasource_test

import (
    "context"
    "errors"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/dogab/vitalstack/api/pkg/datasource"
)

func TestFSVOClient_SearchProducts_Success(t *testing.T) {
    searchFixture := loadTestFixture(t, "testdata/fsvo_search.json")
    foodFixture := loadTestFixture(t, "testdata/fsvo_food.json")

    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        switch {
        case r.URL.Path == "/foods":
            if q := r.URL.Query().Get("search"); q != "broccoli" {
                t.Errorf("expected search=broccoli, got %s", q)
            }
            if lang := r.URL.Query().Get("lang"); lang != "de" {
                t.Errorf("expected lang=de, got %s", lang)
            }
            _, _ = w.Write(searchFixture)
        case strings.HasPrefix(r.URL.Path, "/food/"):
            _, _ = w.Write(foodFixture)
        default:
            t.Errorf("unexpected path: %s", r.URL.Path)
            w.WriteHeader(http.StatusNotFound)
        }
    }))
    defer srv.Close()

    client := datasource.NewFSVOClient(srv.Client(), srv.URL, datasource.WithFSVOLanguage("de"))

    products, err := client.SearchProducts(context.Background(), "broccoli", 10)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(products) != 2 {
        t.Fatalf("expected 2 products, got %d", len(products))
    }

    p := products[0]
    if p.ID != "fsvo-349687" {
        t.Errorf("expected ID fsvo-349687, got %s", p.ID)
    }
    if p.Name != "Broccoli, raw" {
        t.Errorf("expected name 'Broccoli, raw', got %s", p.Name)
    }
    if p.Source != "fsvo" {
        t.Errorf("expected source fsvo, got %s", p.Source)
    }
    if len(p.Categories) != 1 || p.Categories[0] != "Fresh vegetables" {
        t.Errorf("expected category 'Fresh vegetables', got %v", p.Categories)
    }
}

func TestFSVOClient_SearchProducts_NutrientMapping(t *testing.T) {
    searchFixture := loadTestFixture(t, "testdata/fsvo_search.json")
    foodFixture := loadTestFixture(t, "testdata/fsvo_food.json")

    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        if r.URL.Path == "/foods" {
            _, _ = w.Write(searchFixture)
        } else {
            _, _ = w.Write(foodFixture)
        }
    }))
    defer srv.Close()

    client := datasource.NewFSVOClient(srv.Client(), srv.URL, datasource.WithFSVOLanguage("de"))

    products, err := client.SearchProducts(context.Background(), "broccoli", 10)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    macros := products[0].Macros
    if macros.Calories != 39 {
        t.Errorf("expected 39 calories, got %f", macros.Calories)
    }
    if macros.Protein != 3.6 {
        t.Errorf("expected 3.6g protein, got %f", macros.Protein)
    }
    if macros.Carbs != 3.2 {
        t.Errorf("expected 3.2g carbs, got %f", macros.Carbs)
    }
    if macros.Fat != 0.6 {
        t.Errorf("expected 0.6g fat, got %f", macros.Fat)
    }
    if macros.Fiber != 3.2 {
        t.Errorf("expected 3.2g fiber, got %f", macros.Fiber)
    }
}

func TestFSVOClient_LookupBarcode_AlwaysNotFound(t *testing.T) {
    client := datasource.NewFSVOClient(http.DefaultClient, "http://unused", datasource.WithFSVOLanguage("de"))

    _, err := client.LookupBarcode(context.Background(), "0000000000000")
    if !errors.Is(err, datasource.ErrNotFound) {
        t.Fatalf("expected ErrNotFound, got %v", err)
    }
}

func TestFSVOClient_SearchProducts_EmptyResults(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        _, _ = w.Write([]byte(`[]`))
    }))
    defer srv.Close()

    client := datasource.NewFSVOClient(srv.Client(), srv.URL, datasource.WithFSVOLanguage("de"))

    products, err := client.SearchProducts(context.Background(), "zzzznonexistent", 5)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(products) != 0 {
        t.Errorf("expected 0 products, got %d", len(products))
    }
}

func TestFSVOClient_Name(t *testing.T) {
    client := datasource.NewFSVOClient(http.DefaultClient, "http://unused")
    if client.Name() != "fsvo" {
        t.Errorf("expected name 'fsvo', got %s", client.Name())
    }
}
```

### Success Criteria:

#### Automated Verification:
- [x] All tests pass: `cd apps/api-go && go test ./pkg/datasource/ -v`
- [x] Full test suite still passes: `cd apps/api-go && go test ./...`
- [x] Formatting is correct: `cd apps/api-go && gofmt -l .` reports no files

#### Manual Verification:
- [x] None required for this phase — tests are the verification

**Implementation Note**: After all automated verification passes, pause here for manual confirmation before proceeding to documentation.

---

## Phase 3: Architecture Documentation Updates

### Overview
Update the backend architecture documentation to reflect the new four-tier waterfall.

### Changes Required:

#### 1. Architecture Doc Updates
**File**: `apps/api-go/architecture.md`
**Changes**:
1. Update the directory structure comment to add FSVO to `pkg/datasource/` description:
   ```
   │   ├── datasource/            # External API clients (OFF, FSVO, USDA)
   ```
2. Update the Mermaid infrastructure node to include FSVO:
   ```
   Infrastructure[Datasources & Search Engine<br>• Meilisearch Local Cache<br>• Open Food Facts<br>• Swiss FSVO<br>• USDA FoodData Central]
   ```
3. Update the Product Service waterfall description (Section 5):
   ```markdown
   - **Product Service:** Implements a multi-layer waterfall architecture:
     1. **Cache:** Local Meilisearch index (fast, typo-tolerant).
     2. **Primary:** Open Food Facts HTTP API (branded products, barcodes).
     3. **Secondary:** Swiss FSVO food composition database (~1,100 generic foods).
     4. **Tertiary:** USDA FoodData Central HTTP API (450K+ foods).
   ```
4. Add FSVO to the Configuration table:
   ```
   | `fsvo.base-url` | `https://api.webapp...BLV-api` | FSVO API base URL |
   | `fsvo.language` | `de` | FSVO response language (de/en/fr/it) |
   ```

### Success Criteria:

#### Automated Verification:
- [x] Documentation has no broken markdown: Review formatting visually

#### Manual Verification:
- [x] User confirms architecture documents accurately reflect the new implementation

**Implementation Note**: After completing this phase, pause here for final manual confirmation before closing the implementation plan.

---

## Testing Strategy

### Unit Tests:
- `TestFSVOClient_SearchProducts_Success` — two-step search + detail, domain mapping
- `TestFSVOClient_SearchProducts_NutrientMapping` — correct macro extraction by component code
- `TestFSVOClient_LookupBarcode_AlwaysNotFound` — sentinel error returned
- `TestFSVOClient_SearchProducts_EmptyResults` — empty array handled gracefully
- `TestFSVOClient_Name` — returns `"fsvo"`

### Key Edge Cases:
- Empty search results → returns empty slice, no error
- Detail fetch fails for one item → skip that item, return rest
- FSVO API returns non-200 → return meaningful error

### Manual Testing Steps:
1. Start server: `make dev-api`
2. Search generic food: `curl "http://localhost:8080/api/products/search?q=Broccoli"` → expect FSVO-sourced results
3. Barcode lookup: `curl "http://localhost:8080/api/products/barcode/7613035466432"` → expect OFF result (FSVO skipped)
4. Search again → same query returns from Meilisearch cache (faster)

## Performance Considerations

- **Two-step lookup adds latency**: Each FSVO search result requires a detail API call. With the default limit of 10, this means up to 11 HTTP calls (1 search + 10 details). This is acceptable because:
  - FSVO is only hit when Meilisearch and OFF both miss
  - Results get cached in Meilisearch, so repeat searches are fast
  - The FSVO dataset is small (~1,100 foods), so response times are quick
- **Sequential detail fetches**: Could be parallelized with goroutines in the future, but sequential is simpler and FSVO latency is low. Not worth the complexity now.

## Migration Notes

None — no database schema changes. FSVO products are cached in the existing Meilisearch `products` index alongside OFF and USDA entries. The `fsvo-` ID prefix prevents collisions.

## References

- Research: `thoughts/research/research_fsvo-integration.md`
- Similar implementation: [`usda_client.go`](file:///Users/dogab/code/MacroGuard/apps/api-go/pkg/datasource/usda_client.go)
- FSVO API docs: [OpenAPI spec](https://api.webapp.prod.blv.foodcase-services.com/BLV_WebApp_WS/webresources/openapi.json)
