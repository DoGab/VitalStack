package controller

import "github.com/dogab/vitalstack/api/pkg/types"

// BarcodeInput represents a barcode lookup request.
type BarcodeInput struct {
	EAN string `path:"ean" doc:"EAN/UPC barcode" example:"0049000000443"`
}

// SearchProductsInput represents a product search request.
type SearchProductsInput struct {
	Query string `query:"query" required:"true" doc:"Search query" example:"yogurt"`
	Limit int    `query:"limit" default:"10" doc:"Max results to return" example:"10"`
}

// ProductBody represents a product in API responses.
type ProductBody struct {
	ID         string            `json:"id"                    doc:"Product identifier" example:"off-0049000000443"`
	Barcode    string            `json:"barcode"               doc:"EAN/UPC barcode" example:"0049000000443"`
	Name       string            `json:"name"                  doc:"Product name" example:"Caffè Latte"`
	Brand      string            `json:"brand"                 doc:"Brand name" example:"Emmi"`
	ImageURL   string            `json:"image_url,omitempty"   doc:"Product image URL"`
	Source     string            `json:"source"                doc:"Data source" example:"openfoodfacts"`
	NutriScore string            `json:"nutri_score,omitempty" doc:"Nutri-Score grade (A-E)" example:"C"`
	Macros     MacrosPer100gBody `json:"macros"                doc:"Nutritional values per 100g"`
}

// MacrosPer100gBody represents per-100g nutritional data in API responses.
type MacrosPer100gBody struct {
	Calories float64 `json:"calories" doc:"Calories per 100g" example:"56"`
	Protein  float64 `json:"protein"  doc:"Protein per 100g (grams)" example:"3.2"`
	Carbs    float64 `json:"carbs"    doc:"Carbohydrates per 100g (grams)" example:"5.1"`
	Fat      float64 `json:"fat"      doc:"Fat per 100g (grams)" example:"1.8"`
	Fiber    float64 `json:"fiber"    doc:"Fiber per 100g (grams)" example:"0"`
}

// BarcodeOutput is the HTTP response for barcode lookups.
type BarcodeOutput struct {
	Body *ProductBody
}

// SearchProductsOutput is the HTTP response for product search.
type SearchProductsOutput struct {
	Body *SearchProductsOutputBody
}

// SearchProductsOutputBody wraps the list of products with attribution.
type SearchProductsOutputBody struct {
	Products    []ProductBody `json:"products"              doc:"List of matching products"`
	Attribution string        `json:"attribution,omitempty" doc:"Data attribution notice"`
}

// productBodyFromDomain converts a domain Product to an API response body.
func productBodyFromDomain(p types.Product) ProductBody {
	return ProductBody{
		ID:         p.ID,
		Barcode:    p.Barcode,
		Name:       p.Name,
		Brand:      p.Brand,
		ImageURL:   p.ImageURL,
		Source:     p.Source,
		NutriScore: p.NutriScore,
		Macros: MacrosPer100gBody{
			Calories: p.Macros.Calories,
			Protein:  p.Macros.Protein,
			Carbs:    p.Macros.Carbs,
			Fat:      p.Macros.Fat,
			Fiber:    p.Macros.Fiber,
		},
	}
}
