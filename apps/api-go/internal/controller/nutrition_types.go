package controller

// ScanInput represents the scan request body
type ScanInput struct {
	Body *ScanInputBody `json:"body"`
}

type ScanInputBody struct {
	ImageBase64 string `json:"image_base64,omitempty" doc:"Base64 encoded image data"`
}

// MacroData represents nutritional macro information
type MacroData struct {
	Calories int     `json:"calories" example:"450" doc:"Total calories"`
	Protein  float64 `json:"protein" example:"25.5" doc:"Protein in grams"`
	Carbs    float64 `json:"carbs" example:"45.0" doc:"Carbohydrates in grams"`
	Fat      float64 `json:"fat" example:"15.5" doc:"Fat in grams"`
	Fiber    float64 `json:"fiber" example:"5.0" doc:"Fiber in grams"`
}

// ScanOutput represents the scan response
type ScanOutput struct {
	Body *ScanOutputBody `json:"body"`
}

type ScanOutputBody struct {
	FoodName    string     `json:"food_name" example:"Grilled Chicken Salad" doc:"Detected food name"`
	Confidence  float64    `json:"confidence" example:"0.92" doc:"Detection confidence score"`
	Macros      *MacroData `json:"macros" doc:"Nutritional macro information"`
	ServingSize string     `json:"serving_size" example:"1 plate (350g)" doc:"Estimated serving size"`
}
