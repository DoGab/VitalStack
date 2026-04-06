package controller

import "github.com/dogab/vitalstack/api/pkg/types"

// LookupBarcodeInput represents a barcode lookup request.
type LookupBarcodeInput struct {
	Barcode string `path:"barcode" doc:"EAN/UPC barcode to look up" example:"7613035466432"`
	Lang    string `query:"lang" doc:"Preferred language code (e.g. de, fr, en)" example:"de"`
}

// SearchProductsInput represents a product search request.
type SearchProductsInput struct {
	Query string `query:"query" doc:"Free-text search query" example:"nutella"`
	Limit int    `query:"limit" doc:"Maximum number of results" minimum:"1" maximum:"50" default:"10"`
	Lang  string `query:"lang" doc:"Preferred language code (e.g. de, fr, en)" example:"de"`
}

// ProductOutput represents a product in API responses.
type ProductOutput struct {
	ID              string       `json:"id" doc:"Composite product identifier"`
	Barcode         string       `json:"barcode" doc:"EAN/UPC barcode"`
	Name            string       `json:"name" doc:"Product display name"`
	Brand           string       `json:"brand" doc:"Brand name"`
	Categories      []string     `json:"categories" doc:"Product categories"`
	ImageURL        string       `json:"image_url" doc:"Product image URL"`
	Source          string       `json:"source" doc:"Data source identifier"`
	NutriScore      string       `json:"nutri_score" doc:"NutriScore grade (A-E)"`
	ServingSize     string       `json:"serving_size,omitempty" doc:"Human-readable serving size (e.g. 250ml)"`
	ServingQuantity float64      `json:"serving_quantity,omitempty" doc:"Serving quantity in grams/ml"`
	Macros          MacrosOutput `json:"macros" doc:"Nutritional values per 100g"`
}

// MacrosOutput represents per-100g nutritional data in API responses.
type MacrosOutput struct {
	Calories float64 `json:"calories" doc:"Calories per 100g" example:"56"`
	Protein  float64 `json:"protein"  doc:"Protein per 100g (grams)" example:"3.2"`
	Carbs    float64 `json:"carbs"    doc:"Carbohydrates per 100g (grams)" example:"5.1"`
	Fat      float64 `json:"fat"      doc:"Fat per 100g (grams)" example:"1.8"`
	Fiber    float64 `json:"fiber"    doc:"Fiber per 100g (grams)" example:"0"`
}

// BarcodeOutput is the HTTP response for barcode lookups.
type BarcodeOutput struct {
	Body *ProductOutput
}

// SearchProductsOutput is the HTTP response for product search.
type SearchProductsOutput struct {
	Body *SearchProductsOutputBody
}

// SearchProductsOutputBody wraps the list of products with attribution.
type SearchProductsOutputBody struct {
	Products    []ProductOutput `json:"products"              doc:"List of matching products"`
	Attribution string          `json:"attribution,omitempty" doc:"Data attribution notice"`
}

// productOutputFromDomain converts a domain Product to an API response body.
func productOutputFromDomain(p types.Product) ProductOutput {
	return ProductOutput{
		ID:              p.ID,
		Barcode:         p.Barcode,
		Name:            p.Name,
		Brand:           p.Brand,
		Categories:      p.Categories,
		ImageURL:        p.ImageURL,
		Source:          p.Source,
		NutriScore:      p.NutriScore,
		ServingSize:     p.ServingSize,
		ServingQuantity: p.ServingQuantity,
		Macros: MacrosOutput{
			Calories: p.Macros.Calories,
			Protein:  p.Macros.Protein,
			Carbs:    p.Macros.Carbs,
			Fat:      p.Macros.Fat,
			Fiber:    p.Macros.Fiber,
		},
	}
}
