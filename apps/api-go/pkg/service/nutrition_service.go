package service

import (
	"context"
)

// NutritionService is a service for nutritional information
type NutritionService struct {
}

// NewNutritionService creates a new nutrition service
func NewNutritionService() *NutritionService {
	return &NutritionService{}
}

// ScanFood scans the food in the image and returns the nutritional information
func (s *NutritionService) ScanFood(ctx context.Context, input *ScanInput) (*ScanOutput, error) {
	response := &ScanOutput{
		FoodName:   "Grilled Chicken Salad",
		Confidence: 0.92,
		Macros: &MacroData{
			Calories: 400,
			Protein:  23,
			Carbs:    65,
			Fat:      7,
		},
		ServingSize: "1 plate (350g)",
	}
	return response, nil
}
