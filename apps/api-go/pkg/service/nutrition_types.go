package service

type ScanInput struct {
	ImageBase64 string `json:"image_base64,omitempty"`
}

type ScanOutput struct {
	FoodName    string     `json:"food_name"`
	Confidence  float64    `json:"confidence"`
	Macros      *MacroData `json:"macros"`
	ServingSize string     `json:"serving_size"`
}

type MacroData struct {
	Calories int     `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
}
