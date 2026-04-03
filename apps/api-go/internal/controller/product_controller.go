package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/dogab/vitalstack/api/pkg/datasource"
	"github.com/dogab/vitalstack/api/pkg/types"
)

const (
	offAttribution = "Product data sourced from Open Food Facts (openfoodfacts.org) under ODbL."
)

// ProductServicer defines the interface needed by ProductController.
type ProductServicer interface {
	LookupBarcode(ctx context.Context, barcode string) (*types.Product, error)
	SearchProducts(ctx context.Context, query string, limit int) ([]types.Product, error)
}

// ProductController handles product lookup and search HTTP endpoints.
type ProductController struct {
	Service ProductServicer
}

// NewProductController creates a new product controller.
func NewProductController(svc ProductServicer) *ProductController {
	return &ProductController{Service: svc}
}

// Register registers the product endpoints with the Huma API.
func (c *ProductController) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		Path:        "/api/products/barcode/{ean}",
		Method:      http.MethodGet,
		OperationID: "lookup-product-barcode",
		Summary:     "Look up product by barcode",
		Description: "Look up a food product by its EAN/UPC barcode. First checks the local cache, then queries Open Food Facts and USDA FoodData Central.",
		Tags:        []string{"products"},
	}, c.BarcodeHandler)

	huma.Register(api, huma.Operation{
		Path:        "/api/products/search",
		Method:      http.MethodGet,
		OperationID: "search-products",
		Summary:     "Search products",
		Description: "Full-text search across cached and external food product databases. Results are deduplicated and cached for future queries.",
		Tags:        []string{"products"},
	}, c.SearchHandler)
}

// BarcodeHandler handles barcode lookup requests.
func (c *ProductController) BarcodeHandler(ctx context.Context, input *BarcodeInput) (*BarcodeOutput, error) {
	product, err := c.Service.LookupBarcode(ctx, input.EAN)
	if err != nil {
		if errors.Is(err, datasource.ErrNotFound) {
			return nil, huma.Error404NotFound("product not found for barcode: " + input.EAN)
		}
		return nil, huma.Error500InternalServerError("failed to look up product", err)
	}

	body := productBodyFromDomain(*product)
	return &BarcodeOutput{Body: &body}, nil
}

// SearchHandler handles product search requests.
func (c *ProductController) SearchHandler(ctx context.Context, input *SearchProductsInput) (*SearchProductsOutput, error) {
	if input.Limit <= 0 || input.Limit > 50 {
		input.Limit = 10
	}

	products, err := c.Service.SearchProducts(ctx, input.Query, input.Limit)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to search products", err)
	}

	bodies := make([]ProductBody, 0, len(products))
	hasOFF := false
	for _, p := range products {
		bodies = append(bodies, productBodyFromDomain(p))
		if p.Source == "openfoodfacts" {
			hasOFF = true
		}
	}

	resp := &SearchProductsOutputBody{
		Products: bodies,
	}
	if hasOFF {
		resp.Attribution = offAttribution
	}

	return &SearchProductsOutput{Body: resp}, nil
}
