package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dogab/vitalstack/api/pkg/types"
)

const (
	offBaseURL   = "https://ch.openfoodfacts.org"
	offUserAgent = "VitalStack/1.0 (github.com/DoGab/VitalStack)"
	offSource    = "openfoodfacts"
	offFields    = "code,product_name,brands,categories_tags,image_url,nutriments,nutriscore_grade,serving_size,serving_quantity"
)

// offProductResponse represents the Open Food Facts API response for a single product lookup.
type offProductResponse struct {
	Status  int        `json:"status"` // 1 = found, 0 = not found
	Product offProduct `json:"product"`
}

// offSearchResponse represents the Open Food Facts search API response.
type offSearchResponse struct {
	Products []offProduct `json:"products"`
}

// offProduct represents a single product from the Open Food Facts API.
type offProduct struct {
	Code            string        `json:"code"`
	ProductName     string        `json:"product_name"`
	Brands          string        `json:"brands"`
	CategoriesTags  []string      `json:"categories_tags"`
	ImageURL        string        `json:"image_url"`
	NutriScore      string        `json:"nutriscore_grade"`
	Nutriments      offNutriments `json:"nutriments"`
	ServingSize     string        `json:"serving_size"`
	ServingQuantity float64       `json:"serving_quantity"`
}

// offNutriments represents the nutriment data from the Open Food Facts API.
type offNutriments struct {
	EnergyKcal100g    float64 `json:"energy-kcal_100g"`
	Proteins100g      float64 `json:"proteins_100g"`
	Carbohydrates100g float64 `json:"carbohydrates_100g"`
	Fat100g           float64 `json:"fat_100g"`
	Fiber100g         float64 `json:"fiber_100g"`
}

// OFFClient is an HTTP client for the Open Food Facts API.
type OFFClient struct {
	httpClient *http.Client
	baseURL    string
	language   string
	sortBy     string
}

// OFFClientOption defines a functional option for configuring the OFFClient.
type OFFClientOption func(*OFFClient)

// WithBaseURL configures the base URL for the Open Food Facts API.
func WithBaseURL(url string) OFFClientOption {
	return func(c *OFFClient) {
		c.baseURL = url
	}
}

// WithLanguage configures the language code (e.g. "en") for search results.
func WithLanguage(lang string) OFFClientOption {
	return func(c *OFFClient) {
		c.language = lang
	}
}

// WithSortBy configures the sorting method (e.g. "popularity") for search results.
func WithSortBy(sortBy string) OFFClientOption {
	return func(c *OFFClient) {
		c.sortBy = sortBy
	}
}

// NewOFFClient creates a new Open Food Facts client with functional options.
func NewOFFClient(httpClient *http.Client, opts ...OFFClientOption) *OFFClient {
	c := &OFFClient{
		httpClient: httpClient,
		baseURL:    offBaseURL,
		language:   "en",
		sortBy:     "popularity",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SetBaseURL overrides the base URL (used for testing with httptest).
func (c *OFFClient) SetBaseURL(url string) {
	c.baseURL = url
}

// Name returns the datasource identifier.
func (c *OFFClient) Name() string {
	return offSource
}

// LookupBarcode searches for a product by its EAN/UPC barcode.
// The lang parameter overrides the client's default language when non-empty.
func (c *OFFClient) LookupBarcode(ctx context.Context, barcode string, lang string) (*types.Product, error) {
	effectiveLang := c.language
	if lang != "" {
		effectiveLang = lang
	}

	reqURL := fmt.Sprintf("%s/api/v2/product/%s?fields=%s", c.baseURL, url.PathEscape(barcode), offFields)
	if effectiveLang != "" {
		reqURL += "&lc=" + url.QueryEscape(effectiveLang)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("off: creating request: %w", err)
	}
	req.Header.Set("User-Agent", offUserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("off: executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("off: unexpected status %d", resp.StatusCode)
	}

	var result offProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("off: decoding response: %w", err)
	}

	// OFF returns status=0 when product is not found
	if result.Status == 0 {
		return nil, ErrNotFound
	}

	product := offProductToDomain(result.Product)
	return &product, nil
}

// SearchProducts performs a free-text search and returns up to limit results.
// The lang parameter overrides the client's default language when non-empty.
func (c *OFFClient) SearchProducts(ctx context.Context, query string, limit int, lang string) ([]types.Product, error) {
	effectiveLang := c.language
	if lang != "" {
		effectiveLang = lang
	}
	reqURL := fmt.Sprintf(
		"%s/cgi/search.pl?search_terms=%s&json=1&page_size=%d&fields=%s",
		c.baseURL,
		url.QueryEscape(query),
		limit,
		offFields,
	)

	if effectiveLang != "" {
		reqURL += "&lc=" + url.QueryEscape(effectiveLang)
	}
	if c.sortBy != "" {
		reqURL += "&sort_by=" + url.QueryEscape(c.sortBy)
	}

	var resp *http.Response
	var err error

	// Retry logic for 502/503/504 errors which happen when OFF builds a new cache for a search term
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if reqErr != nil {
			return nil, fmt.Errorf("off: creating search request: %w", reqErr)
		}
		req.Header.Set("User-Agent", offUserAgent)

		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("off: executing search request: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			break // Success, exit retry loop
		}

		// Don't leak body on retry
		_ = resp.Body.Close()

		if resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusGatewayTimeout {
			if i < maxRetries-1 {
				// Sleep and retry, OFF might have warmed the cache
				time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
				continue
			}
		}

		return nil, fmt.Errorf("off: search unexpected status %d", resp.StatusCode)
	}
	defer func() { _ = resp.Body.Close() }()

	var result offSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("off: decoding search response: %w", err)
	}

	products := make([]types.Product, 0, len(result.Products))
	for _, p := range result.Products {
		products = append(products, offProductToDomain(p))
	}

	return products, nil
}

// offProductToDomain converts an OFF API product to the domain Product type.
func offProductToDomain(p offProduct) types.Product {
	return types.Product{
		ID:              fmt.Sprintf("off-%s", p.Code),
		Barcode:         p.Code,
		Name:            p.ProductName,
		Brand:           p.Brands,
		Categories:      cleanCategories(p.CategoriesTags),
		ImageURL:        p.ImageURL,
		Source:          offSource,
		NutriScore:      p.NutriScore,
		ServingSize:     p.ServingSize,
		ServingQuantity: p.ServingQuantity,
		Macros: types.MacrosPer100g{
			Calories: p.Nutriments.EnergyKcal100g,
			Protein:  p.Nutriments.Proteins100g,
			Carbs:    p.Nutriments.Carbohydrates100g,
			Fat:      p.Nutriments.Fat100g,
			Fiber:    p.Nutriments.Fiber100g,
		},
	}
}

// cleanCategories strips the "en:" locale prefix from OFF category tags.
func cleanCategories(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}
	cleaned := make([]string, 0, len(tags))
	for _, tag := range tags {
		cleaned = append(cleaned, strings.TrimPrefix(tag, "en:"))
	}
	return cleaned
}
