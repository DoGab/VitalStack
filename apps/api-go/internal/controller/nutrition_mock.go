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
}

// ScanHandler handles the scan request
func (c *NutritionMockController) ScanHandler(ctx context.Context, input *ScanInput) (*ScanOutput, error) {
	return &ScanOutput{
		Body: &ScanOutputBody{
			FoodName:   "Grilled Chicken Salad",
			Confidence: 0.92,
			Macros: &MacroData{
				Calories: 400,
				Protein:  23,
				Carbs:    65,
				Fat:      7,
				Fiber:    5,
			},
			ServingSize: "1 plate (350g)",
		},
	}, nil
}
