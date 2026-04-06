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
	ID            int    `json:"id"` // DBID used for detail lookup
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
	Value     float64       `json:"value"`
	Component fsvoComponent `json:"component"`
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
func (c *FSVOClient) LookupBarcode(_ context.Context, _ string, _ string) (*types.Product, error) {
	return nil, ErrNotFound
}

// SearchProducts searches the FSVO database by name.
// This is a two-step process: search for foods, then fetch detail for each.
// The lang parameter overrides the client's default language when non-empty.
func (c *FSVOClient) SearchProducts(ctx context.Context, query string, limit int, lang string) ([]types.Product, error) {
	effectiveLang := c.language
	if lang != "" {
		effectiveLang = lang
	}

	// Step 1: Search for matching foods (returns metadata only, no macros).
	searchFoods, err := c.searchFoods(ctx, query, limit, effectiveLang)
	if err != nil {
		return nil, err
	}

	// Step 2: Fetch full detail for each search result to get macros.
	products := make([]types.Product, 0, len(searchFoods))
	for _, sf := range searchFoods {
		food, err := c.getFood(ctx, sf.ID, effectiveLang)
		if err != nil {
			// Skip individual failures — don't fail the whole batch.
			continue
		}
		products = append(products, fsvoFoodToDomain(*food))
	}

	return products, nil
}

// searchFoods calls the FSVO /foods search endpoint.
func (c *FSVOClient) searchFoods(ctx context.Context, query string, limit int, lang string) ([]fsvoSearchFood, error) {
	params := url.Values{
		"search": {query},
		"lang":   {lang},
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
func (c *FSVOClient) getFood(ctx context.Context, dbid int, lang string) (*fsvoFood, error) {
	reqURL := fmt.Sprintf("%s/food/%d?lang=%s", c.baseURL, dbid, lang)

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
