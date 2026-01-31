package controller

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dogab/macroguard/api/pkg/service"
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
		Summary:     "Scan food",
		Description: "Scan food",
		Tags:        []string{"nutrition"},
	}, c.ScanHandler)
}

// ScanHandler handles the scan request
func (c *NutritionMockController) ScanHandler(ctx context.Context, input *service.ScanInput) (*service.ScanOutput, error) {
	return nil, nil
}
