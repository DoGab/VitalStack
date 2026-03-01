package controller

// ScanInput represents the scan request body
type ScanInput struct {
	Body *ScanInputBody `json:"body"`
}

type ScanInputBody struct {
	ImageBase64 string  `json:"image_base64" required:"true" doc:"Base64 encoded image data"`
	Description *string `json:"description,omitempty" doc:"Optional meal description for better AI analysis"`
}

// MacroData represents nutritional macro information
type MacroData struct {
	Calories int     `json:"calories" example:"450" doc:"Total calories"`
	Protein  float64 `json:"protein" example:"25.5" doc:"Protein in grams"`
	Carbs    float64 `json:"carbs" example:"45.0" doc:"Carbohydrates in grams"`
	Fat      float64 `json:"fat" example:"15.5" doc:"Fat in grams"`
	Fiber    float64 `json:"fiber" example:"5.0" doc:"Fiber in grams"`
}

// IngredientBody represents a single food component with its nutritional data
type IngredientBody struct {
	Name            string     `json:"name" example:"Grilled Chicken Breast" doc:"Ingredient name"`
	ServingSize     *int       `json:"serving_size,omitempty" example:"150" doc:"Raw serving size generic value"`
	ServingQuantity *float64   `json:"serving_quantity,omitempty" example:"1.5" doc:"Quantity of the serving"`
	ServingUnit     *string    `json:"serving_unit,omitempty" example:"g" doc:"Unit of the serving size (e.g., g, ml)"`
	Macros          *MacroData `json:"macros" doc:"Nutritional macro information for this ingredient"`
}

// ScanOutput represents the scan response
type ScanOutput struct {
	Body *ScanOutputBody `json:"body"`
}

type ScanOutputBody struct {
	FoodName    string           `json:"food_name" example:"Grilled Chicken Salad" doc:"Detected food name"`
	Confidence  float64          `json:"confidence" example:"0.92" doc:"Detection confidence score"`
	Macros      *MacroData       `json:"macros" doc:"Nutritional macro information"`
	ServingSize string           `json:"serving_size" example:"1 plate (350g)" doc:"Estimated serving size"`
	Ingredients []IngredientBody `json:"ingredients" doc:"Breakdown of individual ingredients with their macros"`
}

// LogFoodInput represents the request to log a specific food
type LogFoodInput struct {
	Body *LogFoodInputBody `json:"body"`
}

type LogFoodInputBody struct {
	UserID      *string          `json:"user_id,omitempty" doc:"Optional UUID of the user logging the meal (defaults to auth context if implemented)"`
	FoodName    string           `json:"food_name" example:"Grilled Chicken Salad" doc:"Detected food name"`
	Confidence  float64          `json:"confidence" example:"0.92" doc:"Detection confidence score"`
	Macros      *MacroData       `json:"macros" doc:"Nutritional macro information"`
	Ingredients []IngredientBody `json:"ingredients" doc:"Breakdown of individual ingredients with their macros"`
}

// LogFoodOutput represents the log response
type LogFoodOutput struct {
	Body *LogFoodOutputBody `json:"body"`
}

type LogFoodOutputBody struct {
	Success bool   `json:"success" example:"true" doc:"Indicates whether the log was successfully saved"`
	ID      string `json:"id,omitempty" example:"1" doc:"Database ID of the created food log"`
}

// DailyIntakeInput represents the request to get daily food logs
type DailyIntakeInput struct {
	TzOffset int `query:"tz_offset" default:"0" doc:"Timezone offset in minutes (UTC - Local Time)"`
}

// Meal represents a single logged food item
type Meal struct {
	ID          string           `json:"id" example:"1" doc:"Database ID of the meal"`
	Name        string           `json:"name" example:"Grilled Chicken Salad" doc:"Food name"`
	Time        string           `json:"time" example:"08:15 AM" doc:"Formatted local time of log"`
	Calories    int              `json:"calories" example:"450" doc:"Total calories in the meal"`
	Macros      MacroData        `json:"macros" doc:"Nutritional macro information for the meal"`
	Ingredients []IngredientBody `json:"ingredients,omitempty" doc:"List of ingredients attached to this meal"`
	Confidence  float64          `json:"confidence" example:"0.95" doc:"AI Identification confidence score"`
	ServingSize string           `json:"serving_size" example:"150g" doc:"Total estimated weight or serving size"`
	Emoji       string           `json:"emoji" example:"🥗" doc:"Visual representation emoji"`
	Tag         string           `json:"tag,omitempty" example:"High Protein" doc:"Optional descriptive tag"`
}

// DailyIntakeOutput represents the daily intake response
type DailyIntakeOutput struct {
	Body *DailyIntakeOutputBody `json:"body"`
}

type DailyIntakeOutputBody struct {
	Macros MacroData `json:"macros" doc:"Aggregated nutritional macro information"`
	Meals  []Meal    `json:"meals" doc:"List of meals logged today"`
}

// DeleteLogInput handles removing a scan
type DeleteLogInput struct {
	ID string `path:"id" doc:"ID of the logged meal to delete"`
}

// DeleteLogOutput represents a successful deletion
type DeleteLogOutput struct {
	Body *struct{}
}
