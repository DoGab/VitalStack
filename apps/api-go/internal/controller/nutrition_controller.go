package controller

import (
	"context"
	"net/http"

	"github.com/dogab/macroguard/api/pkg/service"

	"github.com/danielgtaylor/huma/v2"
)

// NutritionServicer is an interface for nutrition services
type NutritionServicer interface {
	ScanFood(ctx context.Context, input *service.ScanInput) (*service.ScanOutput, error)
}

// NutritionController is a controller for nutrition services
type NutritionController struct {
	Service NutritionServicer
}

// NewNutritionController creates a new nutrition controller
func NewNutritionController(service NutritionServicer) *NutritionController {
	return &NutritionController{Service: service}
}

func (c *NutritionController) Register(api huma.API) {
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
func (c *NutritionController) ScanHandler(ctx context.Context, input *ScanInput) (*ScanOutput, error) {
	out := &ScanOutput{}

	req := &service.ScanInput{
		ImageBase64: input.Body.ImageBase64,
	}
	resp, err := c.Service.ScanFood(ctx, req)
	if err != nil {
		return nil, err
	}

	out.Body.FoodName = resp.FoodName
	out.Body.Confidence = resp.Confidence
	out.Body.Macros = &MacroData{
		Calories: resp.Macros.Calories,
		Protein:  resp.Macros.Protein,
		Carbs:    resp.Macros.Carbs,
		Fat:      resp.Macros.Fat,
		Fiber:    resp.Macros.Fiber,
	}
	out.Body.ServingSize = resp.ServingSize

	return out, nil
}
