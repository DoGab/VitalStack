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
	usdaBaseURL = "https://api.nal.usda.gov/fdc/v1"
	usdaSource  = "usda"

	// USDA nutrient IDs for the macros we track.
	nutrientIDEnergy  = 1008
	nutrientIDProtein = 1003
	nutrientIDCarbs   = 1005
	nutrientIDFat     = 1004
	nutrientIDFiber   = 1079
)

// usdaSearchResponse represents the USDA FDC search API response.
type usdaSearchResponse struct {
	Foods []usdaFood `json:"foods"`
}

// usdaFood represents a single food item from the USDA FDC API.
type usdaFood struct {
	FDCId         int            `json:"fdcId"`
	Description   string         `json:"description"`
	BrandOwner    string         `json:"brandOwner"`
	BrandName     string         `json:"brandName"`
	GtinUPC       string         `json:"gtinUpc"`
	FoodCategory  string         `json:"foodCategory"`
	FoodNutrients []usdaNutrient `json:"foodNutrients"`
}

// usdaNutrient represents a single nutrient entry from the USDA FDC API.
type usdaNutrient struct {
	NutrientID   int     `json:"nutrientId"`
	NutrientName string  `json:"nutrientName"`
	Value        float64 `json:"value"`
	UnitName     string  `json:"unitName"`
}

// USDAClient is an HTTP client for the USDA FoodData Central API.
type USDAClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

// NewUSDAClient creates a new USDA FoodData Central client.
func NewUSDAClient(httpClient *http.Client, apiKey string) *USDAClient {
	return &USDAClient{
		httpClient: httpClient,
		apiKey:     apiKey,
		baseURL:    usdaBaseURL,
	}
}

// SetBaseURL overrides the base URL (used for testing with httptest).
func (c *USDAClient) SetBaseURL(url string) {
	c.baseURL = url
}

// Name returns the datasource identifier.
func (c *USDAClient) Name() string {
	return usdaSource
}

// LookupBarcode searches for a branded product by its UPC barcode.
// USDA doesn't have a dedicated barcode endpoint, so we search by UPC
// and filter results to match the exact barcode.
// The lang parameter is accepted to satisfy the interface but is ignored (USDA has no language support).
func (c *USDAClient) LookupBarcode(ctx context.Context, barcode string, _ string) (*types.Product, error) {
	foods, err := c.searchFoods(ctx, barcode, "Branded", 5)
	if err != nil {
		return nil, err
	}

	// Filter results for exact barcode match
	for _, food := range foods {
		if food.GtinUPC == barcode {
			product := usdaFoodToDomain(food)
			return &product, nil
		}
	}

	return nil, ErrNotFound
}

// SearchProducts performs a free-text search and returns up to limit results.
// The lang parameter is accepted to satisfy the interface but is ignored (USDA has no language support).
func (c *USDAClient) SearchProducts(ctx context.Context, query string, limit int, _ string) ([]types.Product, error) {
	foods, err := c.searchFoods(ctx, query, "", limit)
	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0, len(foods))
	for _, food := range foods {
		products = append(products, usdaFoodToDomain(food))
	}

	return products, nil
}

// searchFoods calls the USDA FDC search endpoint.
func (c *USDAClient) searchFoods(ctx context.Context, query, dataType string, limit int) ([]usdaFood, error) {
	params := url.Values{
		"query":    {query},
		"pageSize": {strconv.Itoa(limit)},
		"api_key":  {c.apiKey},
	}
	if dataType != "" {
		params.Set("dataType", dataType)
	}

	reqURL := fmt.Sprintf("%s/foods/search?%s", c.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("usda: creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("usda: executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("usda: unexpected status %d", resp.StatusCode)
	}

	var result usdaSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("usda: decoding response: %w", err)
	}

	return result.Foods, nil
}

// usdaFoodToDomain converts a USDA FDC food item to the domain Product type.
func usdaFoodToDomain(f usdaFood) types.Product {
	brand := f.BrandOwner
	if brand == "" {
		brand = f.BrandName
	}

	var categories []string
	if f.FoodCategory != "" {
		categories = []string{f.FoodCategory}
	}

	return types.Product{
		ID:         fmt.Sprintf("usda-%d", f.FDCId),
		Barcode:    f.GtinUPC,
		Name:       f.Description,
		Brand:      brand,
		Categories: categories,
		Source:     usdaSource,
		Macros:     extractUSDAMacros(f.FoodNutrients),
	}
}

// extractUSDAMacros extracts macronutrient values from USDA nutrient data by matching nutrient IDs.
func extractUSDAMacros(nutrients []usdaNutrient) types.MacrosPer100g {
	var macros types.MacrosPer100g
	for _, n := range nutrients {
		switch n.NutrientID {
		case nutrientIDEnergy:
			macros.Calories = n.Value
		case nutrientIDProtein:
			macros.Protein = n.Value
		case nutrientIDCarbs:
			macros.Carbs = n.Value
		case nutrientIDFat:
			macros.Fat = n.Value
		case nutrientIDFiber:
			macros.Fiber = n.Value
		}
	}
	return macros
}
