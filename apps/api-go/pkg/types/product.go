package types

// Product represents a normalized food product from any datasource.
// Products are identified by a composite ID prefixed with the source name
// (e.g. "off-7613035466432") to prevent collisions across datasources.
type Product struct {
	ID         string        `json:"id"`
	Barcode    string        `json:"barcode"`
	Name       string        `json:"name"`
	Brand      string        `json:"brand"`
	Categories []string      `json:"categories"`
	ImageURL   string        `json:"image_url"`
	Source     string        `json:"source"`      // "openfoodfacts", "fsvo", "usda"
	NutriScore string        `json:"nutri_score"` // A-E (OFF only)
	Macros     MacrosPer100g `json:"macros"`
}

// MacrosPer100g holds nutritional data normalized to per-100g values.
type MacrosPer100g struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
}
