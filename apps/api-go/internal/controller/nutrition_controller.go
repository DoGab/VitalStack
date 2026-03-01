package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dogab/vitalstack/api/internal/middleware"
	"github.com/dogab/vitalstack/api/pkg/service"
)

// NutritionServicer is an interface for nutrition services
type NutritionServicer interface {
	ScanFood(ctx context.Context, input *service.ScanInput) (*service.ScanOutput, error)
	LogFood(ctx context.Context, input *service.LogFoodInput) (*service.LogFoodOutput, error)
	GetDailyIntake(ctx context.Context, userID string, tzOffsetMins int) (*service.DailyIntakeOutput, error)
	DeleteLoggedFood(ctx context.Context, userID string, logID int64) error
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

	huma.Register(api, huma.Operation{
		Path:        "/api/nutrition/log",
		Method:      http.MethodPost,
		OperationID: "log-food",
		Summary:     "Log scanned food",
		Description: "Save approved food analysis to the user's diet log.",
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
func (c *NutritionController) ScanHandler(ctx context.Context, input *ScanInput) (*ScanOutput, error) {
	out := &ScanOutput{Body: &ScanOutputBody{}}

	req := &service.ScanInput{
		ImageBase64: input.Body.ImageBase64,
		Description: input.Body.Description,
	}
	resp, err := c.Service.ScanFood(ctx, req)
	if err != nil {
		return nil, convertServiceErrorToHTTPError(err)
	}

	// Compute totals from ingredients
	totals := resp.TotalMacros()

	out.Body.FoodName = resp.FoodName
	out.Body.Confidence = resp.Confidence
	out.Body.Macros = &MacroData{
		Calories: totals.Calories,
		Protein:  totals.Protein,
		Carbs:    totals.Carbs,
		Fat:      totals.Fat,
		Fiber:    totals.Fiber,
	}
	out.Body.ServingSize = fmt.Sprintf("%dg", resp.TotalWeight())

	// Map ingredients from service to controller type
	out.Body.Ingredients = make([]IngredientBody, len(resp.Ingredients))
	for i, ing := range resp.Ingredients {
		out.Body.Ingredients[i] = IngredientBody{
			Name:            ing.Name,
			ServingSize:     ing.ServingSize,
			ServingQuantity: ing.ServingQuantity,
			ServingUnit:     ing.ServingUnit,
			Macros: &MacroData{
				Calories: ing.Calories,
				Protein:  ing.Protein,
				Carbs:    ing.Carbs,
				Fat:      ing.Fat,
				Fiber:    ing.Fiber,
			},
		}
	}

	return out, nil
}

// LogFoodHandler handles saving an accepted scan to the database
func (c *NutritionController) LogFoodHandler(ctx context.Context, input *LogFoodInput) (*LogFoodOutput, error) {
	// Extract user ID from authenticated context
	uid, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	// Map controller input to service input
	serviceIngredients := make([]service.Ingredient, len(input.Body.Ingredients))
	for i, ing := range input.Body.Ingredients {
		serviceIngredients[i] = service.Ingredient{
			Name:            ing.Name,
			ServingSize:     ing.ServingSize,
			ServingQuantity: ing.ServingQuantity,
			ServingUnit:     ing.ServingUnit,
			Calories:        ing.Macros.Calories,
			Protein:         ing.Macros.Protein,
			Carbs:           ing.Macros.Carbs,
			Fat:             ing.Macros.Fat,
			Fiber:           ing.Macros.Fiber,
		}
	}

	serviceReq := &service.LogFoodInput{
		UserID:     &uid,
		FoodName:   input.Body.FoodName,
		Confidence: input.Body.Confidence,
		Macros: service.MacroData{
			Calories: input.Body.Macros.Calories,
			Protein:  input.Body.Macros.Protein,
			Carbs:    input.Body.Macros.Carbs,
			Fat:      input.Body.Macros.Fat,
			Fiber:    input.Body.Macros.Fiber,
		},
		Ingredients: serviceIngredients,
	}

	resp, err := c.Service.LogFood(ctx, serviceReq)
	if err != nil {
		return nil, convertServiceErrorToHTTPError(err)
	}

	return &LogFoodOutput{
		Body: &LogFoodOutputBody{
			Success: resp.Success,
			ID:      resp.ID,
		},
	}, nil
}

// GetDailyIntakeHandler fetches today's aggregated macros and a list of meals
func (c *NutritionController) GetDailyIntakeHandler(ctx context.Context, input *DailyIntakeInput) (*DailyIntakeOutput, error) {
	// Extract user ID from authenticated context
	userID, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	resp, err := c.Service.GetDailyIntake(ctx, userID, input.TzOffset)
	if err != nil {
		return nil, convertServiceErrorToHTTPError(err)
	}

	// Map Service DTOs to HTTP Output
	meals := make([]Meal, len(resp.Meals))
	for i, m := range resp.Meals {
		// Map service ingredients to controller IngredientBody
		ingBodies := make([]IngredientBody, len(m.Ingredients))
		for j, ing := range m.Ingredients {
			ingBodies[j] = IngredientBody{
				Name:            ing.Name,
				ServingSize:     ing.ServingSize,
				ServingQuantity: ing.ServingQuantity,
				ServingUnit:     ing.ServingUnit,
				Macros: &MacroData{
					Calories: ing.Calories,
					Protein:  ing.Protein,
					Carbs:    ing.Carbs,
					Fat:      ing.Fat,
					Fiber:    ing.Fiber,
				},
			}
		}

		meals[i] = Meal{
			ID:       m.ID,
			Name:     m.Name,
			Time:     m.Time,
			Calories: m.Calories,
			Macros: MacroData{
				Calories: m.Macros.Calories,
				Protein:  m.Macros.Protein,
				Carbs:    m.Macros.Carbs,
				Fat:      m.Macros.Fat,
				Fiber:    m.Macros.Fiber,
			},
			Emoji:       m.Emoji,
			Tag:         m.Tag,
			Ingredients: ingBodies,
		}
	}

	return &DailyIntakeOutput{
		Body: &DailyIntakeOutputBody{
			Macros: MacroData{
				Calories: resp.Macros.Calories,
				Protein:  resp.Macros.Protein,
				Carbs:    resp.Macros.Carbs,
				Fat:      resp.Macros.Fat,
				Fiber:    resp.Macros.Fiber,
			},
			Meals: meals,
		},
	}, nil
}

// DeleteLogHandler extracts the user identity and validates the string ID to int64 for database removal
func (c *NutritionController) DeleteLogHandler(ctx context.Context, input *DeleteLogInput) (*DeleteLogOutput, error) {
	// Extract user ID from authenticated context
	userID, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Unauthorized")
	}

	// Parse int64 from the URL path param
	importStrv, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid meal log ID format")
	}

	err = c.Service.DeleteLoggedFood(ctx, userID, importStrv)
	if err != nil {
		return nil, convertServiceErrorToHTTPError(err)
	}

	return &DeleteLogOutput{}, nil
}
