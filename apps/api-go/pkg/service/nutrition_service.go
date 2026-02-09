package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/dogab/vitalstack/api/pkg/types"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
)

// ErrNotFood is returned when the image does not contain food
var ErrNotFood = errors.New("image does not contain food")

// NutritionService is a service for nutritional information
type NutritionService struct {
	genkit *genkit.Genkit
	flows  map[flowName]*core.Flow[*ScanInput, *ScanOutput, struct{}]
}

// NewNutritionService creates a new nutrition service
func NewNutritionService(genkit *genkit.Genkit) *NutritionService {
	svc := &NutritionService{
		genkit: genkit,
		flows:  map[flowName]*core.Flow[*ScanInput, *ScanOutput, struct{}]{},
	}
	svc.initializeFlows()
	return svc
}

// initializeFlows initializes all AI flows for genkit and stores them in the flows map
func (s *NutritionService) initializeFlows() {
	s.flows = map[flowName]*core.Flow[*ScanInput, *ScanOutput, struct{}]{
		FoodScanFlow: genkit.DefineFlow(s.genkit, string(FoodScanFlow), s.foodScanFlow),
	}
}

// ScanFood scans the food in the image and returns the nutritional information
func (s *NutritionService) ScanFood(ctx context.Context, input *ScanInput) (*ScanOutput, error) {
	slog.Info("received food scan request", "input", input)

	flow := s.flows[FoodScanFlow]
	response, err := flow.Run(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to run food scan flow: %w", err)
	}

	slog.Debug("food scan response", "response", response)

	// Check if the image contains food
	if !response.IsFood {
		return nil, types.NewValidationError(ErrNotFood.Error(), "image_base64", "request.body", "<omitted>")
	}

	return response, nil
}

func (s *NutritionService) foodScanFlow(ctx context.Context, input *ScanInput) (*ScanOutput, error) {
	systemPrompt := `You are an expert nutritionist and food recognition AI.
Analyze the provided image and determine if it contains food.

FIRST: Determine if the image contains food
- Set is_food to true if the image contains any food items
- Set is_food to false if the image does NOT contain food (e.g., objects, people, landscapes, text, documents)
- Set detected_object to describe what you see (e.g., "Grilled Chicken Salad" or "Laptop computer")

IF THE IMAGE CONTAINS FOOD (is_food = true):
Proceed with nutritional analysis:

ANALYSIS STEPS:
1. Identify all visible food items, ingredients, and portion sizes
2. For each ingredient, estimate its weight in grams and calculate individual macros
3. Consider cooking methods (fried, grilled, steamed, etc.) as they affect calories
4. Estimate serving sizes relative to standard references (e.g., a fist ≈ 1 cup, palm ≈ 3oz protein)
5. Sum all ingredient macros to get the total meal macros

OUTPUT REQUIREMENTS:
- is_food: true if image contains food, false otherwise
- detected_object: What you see in the image
- food_name: Overall meal/dish name (empty string if not food)
- confidence: How clearly the food is identifiable (0.0-1.0, or 0.0 if not food)
- serving_size: Total serving in grams or standard units (empty if not food)
- macros: Total combined macros for the entire meal (null if not food)
- ingredients: Array of each component (empty array if not food)

IF THE IMAGE DOES NOT CONTAIN FOOD (is_food = false):
Return minimal response with is_food=false and detected_object describing what you see.

GUIDELINES:
- Always break down complex meals into their visible components
- Use reasonable middle-ground estimates when portions are unclear
- Include fiber in macro calculations when applicable`

	// Build the user prompt text
	userPrompt := "Analyze this image. First determine if it contains food, then provide nutritional information if applicable."
	if input.Description != nil && *input.Description != "" {
		userPrompt = fmt.Sprintf("Analyze this image. Additional context: %s. First determine if it contains food, then provide nutritional information if applicable.", *input.Description)
	}

	// Build the image data URL for multimodal input
	imageDataURL := "data:image/jpeg;base64," + input.ImageBase64

	// Generate structured output using proper multimodal input
	result, _, err := genkit.GenerateData[ScanOutput](ctx, s.genkit,
		ai.WithSystem(systemPrompt),
		ai.WithMessages(
			ai.NewUserMessage(
				ai.NewMediaPart("image/jpeg", imageDataURL),
				ai.NewTextPart(userPrompt),
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze food image: %w", err)
	}

	return result, nil
}
