package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/dogab/vitalstack/api/internal/models"
	"github.com/dogab/vitalstack/api/internal/repository"
	"github.com/dogab/vitalstack/api/pkg/types"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
)

// ErrNotFood is returned when the image does not contain food
var ErrNotFood = errors.New("that doesn't look like food. Please try again with a different image")

// NutritionService is a service for nutritional information
type NutritionService struct {
	genkit      *genkit.Genkit
	flows       map[flowName]*core.Flow[*ScanInput, *ScanOutput, struct{}]
	foodLogRepo repository.FoodLogRepository
	mockScan    bool // If true, ScanFood returns diverse dummy data
}

// NutritionServiceOption defines a functional option for configuring the service
type NutritionServiceOption func(*NutritionService)

// WithMockScan sets whether the Genkit flow should be bypassed to return mock meal responses
func WithMockScan(mock bool) NutritionServiceOption {
	return func(s *NutritionService) {
		s.mockScan = mock
	}
}

// NewNutritionService creates a new nutrition service
func NewNutritionService(genkit *genkit.Genkit, foodLogRepo repository.FoodLogRepository, opts ...NutritionServiceOption) *NutritionService {
	svc := &NutritionService{
		genkit:      genkit,
		flows:       map[flowName]*core.Flow[*ScanInput, *ScanOutput, struct{}]{},
		foodLogRepo: foodLogRepo,
	}

	for _, opt := range opts {
		opt(svc)
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

	if s.mockScan {
		slog.Info("returning dynamically mocked scan data (bypassing Genkit)")
		return s.generateMockScan(), nil
	}

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
Identify each visible ingredient and estimate its individual macros.

ANALYSIS STEPS:
1. Identify all visible food items and ingredients
2. For each ingredient, estimate its weight in grams
3. For each ingredient, calculate its individual macros (calories, protein, carbs, fat, fiber)
4. Consider cooking methods (fried, grilled, steamed) — reflect them in the ingredient macros
5. Estimate portion sizes relative to standard references (e.g., a fist ≈ 1 cup, palm ≈ 3oz protein)

OUTPUT REQUIREMENTS:
- is_food: true if image contains food, false otherwise
- detected_object: What you see in the image
- food_name: Overall meal/dish name (empty string if not food)
- confidence: How clearly the food is identifiable (0.0-1.0, or 0.0 if not food)
- ingredients: Array of each component with:
  - name: Ingredient name (e.g., "Grilled Chicken Breast")
  - serving_size: Integer representing the raw unit size or standardized amount (e.g. 100)
  - serving_quantity: Float representing how many servings are present (e.g., 1.5)
  - serving_unit: String representing the unit (e.g., "g", "slice", "cup", "oz")
  - calories: Calories for this ingredient at the estimated weight
  - protein: Protein in grams
  - carbs: Carbohydrates in grams
  - fat: Fat in grams
  - fiber: Fiber in grams

IMPORTANT: Do NOT return total macros. Only return per-ingredient data.
Total macros will be computed by summing all ingredients.

IF THE IMAGE DOES NOT CONTAIN FOOD (is_food = false):
Return minimal response with is_food=false and detected_object describing what you see.

GUIDELINES:
- Always break down complex meals into their visible components
- Use reasonable middle-ground estimates when portions are unclear
- Include cooking oils, sauces, and dressings as separate ingredients when visible
- Include fiber in macro calculations when applicable`

	// Build the user prompt text
	userPrompt := "Analyze this image. First determine if it contains food, then identify each ingredient with its macros."
	if input.Description != nil && *input.Description != "" {
		userPrompt = fmt.Sprintf("Analyze this image. Additional context: %s. First determine if it contains food, then identify each ingredient with its macros.", *input.Description)
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

// LogFood handles saving an accepted scan to the database
func (s *NutritionService) LogFood(ctx context.Context, input *LogFoodInput) (*LogFoodOutput, error) {
	if s.foodLogRepo == nil {
		return nil, errors.New("database repository is not configured")
	}

	// Create the FoodLog database model
	dbLog := &models.FoodLog{
		FoodName:            input.FoodName,
		DetectionConfidence: input.Confidence,
		UserID:              input.UserID,
		Calories:            input.Macros.Calories,
		Protein:             input.Macros.Protein,
		Carbs:               input.Macros.Carbs,
		Fat:                 input.Macros.Fat,
		Fiber:               input.Macros.Fiber,
		CreatedAt:           time.Now().UTC(),
	}

	// Bundle the ingredients directly into a slice
	var dbIngredients []models.FoodLogIngredient
	for _, ing := range input.Ingredients {
		dbIngredients = append(dbIngredients, models.FoodLogIngredient{
			Name:            ing.Name,
			ServingSize:     ing.ServingSize,
			ServingQuantity: ing.ServingQuantity,
			ServingUnit:     ing.ServingUnit,
			Calories:        ing.Calories,
			Protein:         ing.Protein,
			Carbs:           ing.Carbs,
			Fat:             ing.Fat,
			Fiber:           ing.Fiber,
		})
	}

	newLogID, err := s.foodLogRepo.CreateFoodLogWithIngredients(ctx, dbLog, dbIngredients)
	if err != nil {
		slog.Error("Failed to save food log atomically to Supabase", "error", err)
		return nil, fmt.Errorf("failed to save food log: %w", err)
	}

	return &LogFoodOutput{
		Success: true,
		ID:      strconv.FormatInt(newLogID, 10),
	}, nil
}

// GetDailyIntake retrieves aggregated daily macro data and a list of meals
func (s *NutritionService) GetDailyIntake(ctx context.Context, userID string, tzOffsetMins int, targetDateStr string) (*DailyIntakeOutput, error) {
	if s.foodLogRepo == nil {
		return nil, errors.New("database repository is not configured")
	}

	// Calculate bounds based on tzOffset
	// tzOffsetMins is the number of minutes to add to UTC to get local time (e.g. +120 for GMT+2).
	// Wait, timezoneOffset in JS is UTC - Local in minutes.
	// We'll assume the frontend sends the timezone offset in minutes: localTime = UTC - tzOffsetMins
	nowUTC := time.Now().UTC()
	// To get local time of the user:
	localTimeOrigin := nowUTC.Add(time.Duration(-tzOffsetMins) * time.Minute)

	if targetDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", targetDateStr)
		if err == nil {
			// Set the local time origin to the requested date (keeping the time at 00:00:00)
			localTimeOrigin = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
		}
	}

	// Start of local day
	localStartOfDay := time.Date(localTimeOrigin.Year(), localTimeOrigin.Month(), localTimeOrigin.Day(), 0, 0, 0, 0, time.UTC)
	// End of local day
	localEndOfDay := localStartOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	// Shift bounds back to UTC for database querying
	utcStartOfDay := localStartOfDay.Add(time.Duration(tzOffsetMins) * time.Minute)
	utcEndOfDay := localEndOfDay.Add(time.Duration(tzOffsetMins) * time.Minute)

	logs, err := s.foodLogRepo.GetDailyFoodLogs(ctx, userID, utcStartOfDay, utcEndOfDay)
	if err != nil {
		slog.Error("Failed to get daily food logs from Supabase", "error", err)
		return nil, fmt.Errorf("failed to get daily food logs: %w", err)
	}

	var totalMacros MacroData
	var meals []Meal

	for _, log := range logs {
		totalMacros.Calories += log.Calories
		totalMacros.Protein += log.Protein
		totalMacros.Carbs += log.Carbs
		totalMacros.Fat += log.Fat
		totalMacros.Fiber += log.Fiber

		// Map deeply nested DB ingredients back to business logic ingredients
		var mappedIngredients []Ingredient
		for _, ing := range log.Ingredients {
			mappedIngredients = append(mappedIngredients, Ingredient{
				Name:            ing.Name,
				ServingSize:     ing.ServingSize,
				ServingQuantity: ing.ServingQuantity,
				ServingUnit:     ing.ServingUnit,
				Calories:        ing.Calories,
				Protein:         ing.Protein,
				Carbs:           ing.Carbs,
				Fat:             ing.Fat,
				Fiber:           ing.Fiber,
			})
		}

		// Calculate local time for the meal display
		mealLocalTime := log.CreatedAt.UTC().Add(time.Duration(-tzOffsetMins) * time.Minute)

		dummyScan := &ScanOutput{Ingredients: mappedIngredients}
		calcWeight := dummyScan.TotalWeight()
		servingStr := "Unknown amount"
		if calcWeight > 0 {
			servingStr = fmt.Sprintf("%dg", calcWeight)
		}

		meals = append(meals, Meal{
			ID:          strconv.FormatInt(log.ID, 10),
			Name:        log.FoodName,
			Time:        mealLocalTime.Format("03:04 PM"),
			Calories:    log.Calories,
			Confidence:  log.DetectionConfidence,
			ServingSize: servingStr,
			Macros: MacroData{
				Calories: log.Calories,
				Protein:  log.Protein,
				Carbs:    log.Carbs,
				Fat:      log.Fat,
				Fiber:    log.Fiber,
			},
			Ingredients: mappedIngredients,
			Emoji:       "🍽️",
			Tag:         "",
		})
	}

	// Reverse meals array so newest is first (since we removed database ordering)
	for i, j := 0, len(meals)-1; i < j; i, j = i+1, j-1 {
		meals[i], meals[j] = meals[j], meals[i]
	}

	return &DailyIntakeOutput{
		Macros: totalMacros,
		Meals:  meals,
	}, nil
}

// DeleteLoggedFood deletes a logged food entry by ID
func (s *NutritionService) DeleteLoggedFood(ctx context.Context, userID string, logID int64) error {
	if s.foodLogRepo == nil {
		return errors.New("database repository is not configured")
	}

	err := s.foodLogRepo.DeleteFoodLog(ctx, userID, logID)
	if err != nil {
		slog.Error("Failed to delete food log from Supabase", "error", err, "logID", logID)
		return fmt.Errorf("failed to delete food log: %w", err)
	}

	slog.Info("Successfully deleted food log", "logID", logID, "userID", userID)
	return nil
}

// GetHistory retrieves aggregated daily macro data over a specified number of past days
func (s *NutritionService) GetHistory(ctx context.Context, userID string, tzOffsetMins int, days int) (*HistoryOutput, error) {
	if s.foodLogRepo == nil {
		return nil, errors.New("database repository is not configured")
	}

	nowUTC := time.Now().UTC()
	localTimeOrigin := nowUTC.Add(time.Duration(-tzOffsetMins) * time.Minute)

	// End is the end of today
	localEndOfToday := time.Date(localTimeOrigin.Year(), localTimeOrigin.Month(), localTimeOrigin.Day(), 23, 59, 59, 999999999, time.UTC)
	// Start is X days ago (including today as one of the days)
	localStartOfRange := time.Date(localTimeOrigin.Year(), localTimeOrigin.Month(), localTimeOrigin.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -(days - 1))

	utcStartOfRange := localStartOfRange.Add(time.Duration(tzOffsetMins) * time.Minute)
	utcEndOfRange := localEndOfToday.Add(time.Duration(tzOffsetMins) * time.Minute)

	logs, err := s.foodLogRepo.GetDailyFoodLogs(ctx, userID, utcStartOfRange, utcEndOfRange)
	if err != nil {
		slog.Error("Failed to get historical food logs from Supabase", "error", err)
		return nil, fmt.Errorf("failed to get historical food logs: %w", err)
	}

	// Initialize days map/slice
	dailySummariesMap := make(map[string]*DailySummary)

	// Pre-fill last X days to ensure we have all days even if no logs exist
	var orderedDates []string
	for i := days - 1; i >= 0; i-- {
		d := localEndOfToday.AddDate(0, 0, -i).Format("2006-01-02")
		orderedDates = append(orderedDates, d)
		dailySummariesMap[d] = &DailySummary{
			Date:   d,
			Macros: MacroData{},
		}
	}

	var overallTotals MacroData

	for _, log := range logs {
		// Calculate local time for the meal to put it in the correct bucket
		mealLocalTime := log.CreatedAt.UTC().Add(time.Duration(-tzOffsetMins) * time.Minute)
		dateStr := mealLocalTime.Format("2006-01-02")

		if summary, exists := dailySummariesMap[dateStr]; exists {
			summary.Macros.Calories += log.Calories
			summary.Macros.Protein += log.Protein
			summary.Macros.Carbs += log.Carbs
			summary.Macros.Fat += log.Fat
			summary.Macros.Fiber += log.Fiber
		}

		overallTotals.Calories += log.Calories
		overallTotals.Protein += log.Protein
		overallTotals.Carbs += log.Carbs
		overallTotals.Fat += log.Fat
		overallTotals.Fiber += log.Fiber
	}

	// Compute averages
	var averages MacroData
	if days > 0 {
		averages.Calories = int(math.Round(float64(overallTotals.Calories) / float64(days)))
		averages.Protein = overallTotals.Protein / float64(days)
		averages.Carbs = overallTotals.Carbs / float64(days)
		averages.Fat = overallTotals.Fat / float64(days)
		averages.Fiber = overallTotals.Fiber / float64(days)
	}

	// Format averages
	averages.Protein = math.Round(averages.Protein*10) / 10
	averages.Carbs = math.Round(averages.Carbs*10) / 10
	averages.Fat = math.Round(averages.Fat*10) / 10
	averages.Fiber = math.Round(averages.Fiber*10) / 10

	var daysList []DailySummary
	for _, dateStr := range orderedDates {
		summary := dailySummariesMap[dateStr]
		// Round the daily macros
		summary.Macros.Protein = math.Round(summary.Macros.Protein*10) / 10
		summary.Macros.Carbs = math.Round(summary.Macros.Carbs*10) / 10
		summary.Macros.Fat = math.Round(summary.Macros.Fat*10) / 10
		summary.Macros.Fiber = math.Round(summary.Macros.Fiber*10) / 10
		daysList = append(daysList, *summary)
	}

	return &HistoryOutput{
		Averages: averages,
		Days:     daysList,
	}, nil
}

// ptr is a helper function to create pointers to values
func ptr[T any](v T) *T {
	return &v
}

// generateMockScan returns a randomized mock scan output from 3 predefined meals
func (s *NutritionService) generateMockScan() *ScanOutput {
	// 3 hardcoded diverse mock meals
	meals := []ScanOutput{
		{
			IsFood:         true,
			DetectedObject: "A salad bowl with grilled chicken",
			FoodName:       "Grilled Chicken Salad",
			Confidence:     0.95,
			Ingredients: []Ingredient{
				{Name: "Grilled Chicken Breast", Calories: 248, Protein: 38, Carbs: 0, Fat: 10, Fiber: 0, ServingSize: ptr(150), ServingQuantity: ptr(1.0), ServingUnit: ptr("g")},
				{Name: "Mixed Greens", Calories: 20, Protein: 2, Carbs: 3, Fat: 0, Fiber: 2, ServingSize: ptr(100), ServingQuantity: ptr(1.0), ServingUnit: ptr("g")},
				{Name: "Cherry Tomatoes", Calories: 18, Protein: 1, Carbs: 4, Fat: 0, Fiber: 1, ServingSize: ptr(100), ServingQuantity: ptr(0.5), ServingUnit: ptr("cup")},
				{Name: "Feta Cheese", Calories: 105, Protein: 6, Carbs: 2, Fat: 8, Fiber: 0, ServingSize: ptr(28), ServingQuantity: ptr(1.5), ServingUnit: ptr("g")},
				{Name: "Olive Oil Dressing", Calories: 80, Protein: 0, Carbs: 1, Fat: 9, Fiber: 0, ServingSize: ptr(15), ServingQuantity: ptr(1.0), ServingUnit: ptr("ml")},
				{Name: "Cucumber", Calories: 5, Protein: 0, Carbs: 1, Fat: 0, Fiber: 0, ServingSize: ptr(50), ServingQuantity: ptr(0.5), ServingUnit: ptr("cup")},
			},
		},
		{
			IsFood:         true,
			DetectedObject: "A bowl of oatmeal topped with fresh berries",
			FoodName:       "Blueberry Oatmeal",
			Confidence:     0.92,
			Ingredients: []Ingredient{
				{Name: "Rolled Oats", Calories: 150, Protein: 5, Carbs: 27, Fat: 3, Fiber: 4, ServingSize: ptr(40), ServingQuantity: ptr(1.0), ServingUnit: ptr("g")},
				{Name: "Almond Milk", Calories: 30, Protein: 1, Carbs: 1, Fat: 2.5, Fiber: 0, ServingSize: ptr(240), ServingQuantity: ptr(0.5), ServingUnit: ptr("ml")},
				{Name: "Blueberries", Calories: 42, Protein: 0.5, Carbs: 11, Fat: 0.2, Fiber: 1.8, ServingSize: ptr(74), ServingQuantity: ptr(1.0), ServingUnit: ptr("cup")},
				{Name: "Chia Seeds", Calories: 60, Protein: 2, Carbs: 5, Fat: 4, Fiber: 4, ServingSize: ptr(15), ServingQuantity: ptr(1.0), ServingUnit: ptr("tbsp")},
				{Name: "Honey", Calories: 64, Protein: 0, Carbs: 17, Fat: 0, Fiber: 0, ServingSize: ptr(21), ServingQuantity: ptr(1.0), ServingUnit: ptr("tbsp")},
			},
		},
		{
			IsFood:         true,
			DetectedObject: "A bowl of rice with pan-seared salmon",
			FoodName:       "Salmon Rice Bowl",
			Confidence:     0.97,
			Ingredients: []Ingredient{
				{Name: "Seared Salmon", Calories: 280, Protein: 25, Carbs: 0, Fat: 18, Fiber: 0, ServingSize: ptr(140), ServingQuantity: ptr(1.0), ServingUnit: ptr("g")},
				{Name: "Jasmine Rice", Calories: 205, Protein: 4, Carbs: 45, Fat: 0.4, Fiber: 0.6, ServingSize: ptr(158), ServingQuantity: ptr(1.0), ServingUnit: ptr("cup")},
				{Name: "Avocado", Calories: 160, Protein: 2, Carbs: 8, Fat: 15, Fiber: 6, ServingSize: ptr(100), ServingQuantity: ptr(0.5), ServingUnit: ptr("medium")},
				{Name: "Sesame Seeds", Calories: 52, Protein: 1.6, Carbs: 2.1, Fat: 4.5, Fiber: 1.1, ServingSize: ptr(9), ServingQuantity: ptr(1.0), ServingUnit: ptr("tbsp")},
				{Name: "Soy Sauce", Calories: 9, Protein: 1.3, Carbs: 0.8, Fat: 0, Fiber: 0, ServingSize: ptr(15), ServingQuantity: ptr(1.0), ServingUnit: ptr("tbsp")},
			},
		},
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(meals))))
	return &meals[n.Int64()]
}
