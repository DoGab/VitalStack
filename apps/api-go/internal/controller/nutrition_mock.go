package controller

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// NutritionMockController is a controller for mocking the nutrition handlers
type NutritionMockController struct{}

// NewNutritionMockController creates a new nutrition mock controller
func NewNutritionMockController() *NutritionMockController {
	return &NutritionMockController{}
}

// Register registers the nutrition mock controller with the API
func (c *NutritionMockController) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		Path:        "/api/nutrition/scan",
		Method:      http.MethodPost,
		OperationID: "scan-food",
		Summary:     "Scan food image for nutritional information",
		Description: "Upload a base64-encoded food image and optionally provide a description. Returns detected food name and macro breakdown.",
		Tags:        []string{"nutrition"},
	}, c.ScanHandler)

	huma.Register(api, huma.Operation{
		Path:        "/api/nutrition/log",
		Method:      http.MethodPost,
		OperationID: "log-food",
		Summary:     "Log scanned food",
		Description: "Log a previously scanned food item to the database",
		Tags:        []string{"nutrition"},
	}, c.LogFoodHandler)

	huma.Register(api, huma.Operation{
		Path:        "/api/nutrition/daily",
		Method:      http.MethodGet,
		OperationID: "get-daily-intake",
		Summary:     "Get daily intake",
		Description: "Fetch the user's aggregated daily macros and logged meals for today.",
		Tags:        []string{"nutrition"},
	}, c.GetDailyIntakeHandler)

	huma.Register(api, huma.Operation{
		Path:        "/api/nutrition/log/{id}",
		Method:      http.MethodDelete,
		OperationID: "delete-log",
		Summary:     "Delete meal log",
		Description: "Permanently removes a meal and its scanned ingredients from the user's diary.",
		Tags:        []string{"nutrition"},
	}, c.DeleteLogHandler)
}

// ScanHandler handles the scan request
func (c *NutritionMockController) ScanHandler(ctx context.Context, input *ScanInput) (*ScanOutput, error) {
	return &ScanOutput{
		Body: &ScanOutputBody{
			FoodName:   "Grilled Chicken Salad",
			Confidence: 0.95,
			Macros: &MacroData{
				Calories: 476,
				Protein:  47,
				Carbs:    11,
				Fat:      27,
				Fiber:    3,
			},
			ServingSize: "400g",
			Ingredients: []IngredientBody{
				{
					Name:            "Grilled Chicken Breast",
					ServingSize:     ptr(150),
					ServingQuantity: ptr(1.0),
					ServingUnit:     ptr("g"),
					Macros: &MacroData{
						Calories: 248,
						Protein:  38,
						Carbs:    0,
						Fat:      10,
						Fiber:    0,
					},
				},
				{
					Name:            "Mixed Greens",
					ServingSize:     ptr(100),
					ServingQuantity: ptr(1.0),
					ServingUnit:     ptr("g"),
					Macros: &MacroData{
						Calories: 20,
						Protein:  2,
						Carbs:    3,
						Fat:      0,
						Fiber:    2,
					},
				},
				{
					Name:            "Cherry Tomatoes",
					ServingSize:     ptr(100),
					ServingQuantity: ptr(0.5),
					ServingUnit:     ptr("cup"),
					Macros: &MacroData{
						Calories: 18,
						Protein:  1,
						Carbs:    4,
						Fat:      0,
						Fiber:    1,
					},
				},
				{
					Name:            "Feta Cheese",
					ServingSize:     ptr(28),
					ServingQuantity: ptr(1.5),
					ServingUnit:     ptr("g"),
					Macros: &MacroData{
						Calories: 105,
						Protein:  6,
						Carbs:    2,
						Fat:      8,
						Fiber:    0,
					},
				},
				{
					Name:            "Olive Oil Dressing",
					ServingSize:     ptr(15),
					ServingQuantity: ptr(1.0),
					ServingUnit:     ptr("ml"),
					Macros: &MacroData{
						Calories: 80,
						Protein:  0,
						Carbs:    1,
						Fat:      9,
						Fiber:    0,
					},
				},
				{
					Name:            "Cucumber",
					ServingSize:     ptr(50),
					ServingQuantity: ptr(0.5),
					ServingUnit:     ptr("cup"),
					Macros: &MacroData{
						Calories: 5,
						Protein:  0,
						Carbs:    1,
						Fat:      0,
						Fiber:    0,
					},
				},
			},
		},
	}, nil
}

// LogFoodHandler handles the food logging request
func (c *NutritionMockController) LogFoodHandler(ctx context.Context, input *LogFoodInput) (*LogFoodOutput, error) {
	return &LogFoodOutput{
		Body: &LogFoodOutputBody{
			Success: true,
			ID:      "mock-id-1234",
		},
	}, nil
}

// GetDailyIntakeHandler returns mock daily intake
func (c *NutritionMockController) GetDailyIntakeHandler(ctx context.Context, input *DailyIntakeInput) (*DailyIntakeOutput, error) {
	out := &DailyIntakeOutput{
		Body: &DailyIntakeOutputBody{
			Macros: MacroData{
				Calories: 840,
				Protein:  65,
				Carbs:    90,
				Fat:      45,
				Fiber:    12,
			},
			Meals: []Meal{
				{
					ID:          "log-1234",
					Name:        "Grilled Chicken Salad",
					Time:        "12:30 PM",
					Calories:    480,
					Confidence:  0.95,
					ServingSize: "350g",
					Emoji:       "🥗",
					Tag:         "High Protein",
					Ingredients: []IngredientBody{
						{
							Name:            "Grilled Chicken Breast",
							ServingSize:     ptr(150),
							ServingQuantity: ptr(1.0),
							ServingUnit:     ptr("g"),
							Macros: &MacroData{
								Calories: 248,
								Protein:  38,
								Carbs:    0,
								Fat:      10,
								Fiber:    0,
							},
						},
						{
							Name:            "Mixed Greens",
							ServingSize:     ptr(100),
							ServingQuantity: ptr(1.0),
							ServingUnit:     ptr("g"),
							Macros: &MacroData{
								Calories: 20,
								Protein:  2,
								Carbs:    3,
								Fat:      0,
								Fiber:    2,
							},
						},
					},
				},
				{
					ID:          "log-1235",
					Name:        "Morning Berry Smoothie",
					Time:        "08:15 AM",
					Calories:    360,
					Confidence:  0.9,
					ServingSize: "450ml",
					Emoji:       "🫐",
					Tag:         "Antioxidants",
				},
			},
		},
	}

	return out, nil
}

// DeleteLogHandler handles the mock deletion request
func (c *NutritionMockController) DeleteLogHandler(ctx context.Context, input *DeleteLogInput) (*DeleteLogOutput, error) {
	return &DeleteLogOutput{}, nil
}

// ptr is a helper function to create pointers to values
func ptr[T any](v T) *T {
	return &v
}
