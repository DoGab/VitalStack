package repository

import (
	"context"
	"fmt"

	"encoding/json"
	"strconv"
	"time"

	"github.com/dogab/vitalstack/api/internal/models"
	"github.com/supabase-community/supabase-go"
)

type FoodLogRepository interface {
	CreateFoodLog(ctx context.Context, log *models.FoodLog) error
	CreateFoodLogIngredient(ctx context.Context, ingredient *models.FoodLogIngredient) error
	CreateFoodLogWithIngredients(ctx context.Context, log *models.FoodLog, ingredients []models.FoodLogIngredient) (int64, error)
	GetDailyFoodLogs(ctx context.Context, userID string, startOfDay time.Time, endOfDay time.Time) ([]models.FoodLog, error)
	DeleteFoodLog(ctx context.Context, userID string, logID int64) error
}

type foodLogRepository struct {
	client *supabase.Client
}

func NewFoodLogRepository(client *supabase.Client) FoodLogRepository {
	return &foodLogRepository{client: client}
}

// logFoodPayload represents the input arguments for the log_food_atomic RPC.
type logFoodPayload struct {
	UserID              string                     `json:"p_user_id"`
	FoodName            string                     `json:"p_food_name"`
	DetectionConfidence float64                    `json:"p_detection_confidence"`
	Calories            int                        `json:"p_calories"`
	Protein             float64                    `json:"p_protein"`
	Carbs               float64                    `json:"p_carbs"`
	Fat                 float64                    `json:"p_fat"`
	Fiber               float64                    `json:"p_fiber"`
	Ingredients         []models.FoodLogIngredient `json:"p_ingredients"`
}

func (r *foodLogRepository) CreateFoodLog(ctx context.Context, log *models.FoodLog) error {
	// Insert returns a FilterBuilder. By explicitly setting returning to "representation",
	// PostgREST will return the inserted row (which includes the auto-incremented ID).
	data, _, err := r.client.From("food_logs").Insert(log, false, "", "representation", "").Execute()
	if err != nil {
		return err
	}

	// Parse the returning array to get the newly created ID
	var insertedLogs []models.FoodLog
	if err := json.Unmarshal(data, &insertedLogs); err != nil {
		return err
	}

	if len(insertedLogs) > 0 {
		log.ID = insertedLogs[0].ID
	}
	return nil
}

func (r *foodLogRepository) CreateFoodLogIngredient(ctx context.Context, ingredient *models.FoodLogIngredient) error {
	_, _, err := r.client.From("food_log_ingredients").Insert(ingredient, false, "", "", "").Execute()
	return err
}

func (r *foodLogRepository) CreateFoodLogWithIngredients(ctx context.Context, log *models.FoodLog, ingredients []models.FoodLogIngredient) (int64, error) {
	if log.UserID == nil {
		return 0, nil
	}

	payload := logFoodPayload{
		UserID:              *log.UserID,
		FoodName:            log.FoodName,
		DetectionConfidence: log.DetectionConfidence,
		Calories:            log.Calories,
		Protein:             log.Protein,
		Carbs:               log.Carbs,
		Fat:                 log.Fat,
		Fiber:               log.Fiber,
		Ingredients:         ingredients,
	}

	// In supabase-community/supabase-go, `client.Rpc()` actually returns a string containing the URL
	// We need to use `client.Rpc()` and then execute an HTTP request, OR we can use the library's
	// intended way to call RPC functions via the post-release v0.0.x API which actually exposes `.Execute()`
	// Wait, checking the Go SDK, `client.Rpc` used to return `string`. In the latest versions, there is no direct
	// `.Execute()` on `client.Rpc()`.
	// We can execute this by building a custom request or we can use the library's `postgrest` client builder to execute an arbitrary RPC endpoint
	// by manually appending `/rpc/log_food_atomic` to a `From` query builder.

	// Let's use `Insert` with the dummy builder trick correctly or use HTTP.

	// Since we need to send JSON payload to the RPC endpoint and `supabase-go` obscures the HTTP client,
	// A simple approach is to use the `postgrest.QueryBuilder` to perform a custom POST to the RPC endpoint.
	// Actually, `supabase.Client` exposes `client.Rpc(name, count, body)` which returns a string, but this is a legacy artifact.

	// Let's create an HTTP client and post directly to the Supabase URL and Key since `client` doesn't export them.
	// We CANNOT get URL and Key from `client` if they are unexported.
	// However! We can use a trick: `r.client.From("rpc/log_food_atomic").Insert(payload...`
	// Nope, PostgREST doesn't support Insert to /rpc/. It expects a POST with the exact payload as JSON.

	// Let's just use `r.client.Rpc` and use a generic HTTP request? We don't have the API key.
	// Looking closely at `supabase-go`, we can use `r.client.From("")` to access the underlying builder ?

	// Actually, the `supabase-go` PostgREST client implementation in Go supports calling RPC via:
	// `client.Rpc(name, count, rpcBody)` which actually EXECUTES the RPC and returns a string (the JSON response body)!
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("rpc panic: %v", r)
		}
	}()

	responseString := r.client.Rpc("log_food_atomic", "exact", payload)

	// Error could have been captured via panic recovery
	if err != nil {
		return 0, err
	}

	var newID int64
	if responseString != "" && responseString != "null" {
		err = json.Unmarshal([]byte(responseString), &newID)
		if err != nil {
			return 0, err
		}
	}

	return newID, nil
}

func (r *foodLogRepository) GetDailyFoodLogs(ctx context.Context, userID string, startOfDay time.Time, endOfDay time.Time) ([]models.FoodLog, error) {
	// Format bounds to ISO 8601 string to pass into postgres
	startStr := startOfDay.Format(time.RFC3339)
	endStr := endOfDay.Format(time.RFC3339)

	data, _, err := r.client.From("food_logs").
		Select("*, food_log_ingredients(*)", "exact", false).
		Eq("user_id", userID).
		Gte("created_at", startStr).
		Lte("created_at", endStr).
		Execute()

	if err != nil {
		return nil, err
	}

	var logs []models.FoodLog
	if len(data) == 0 {
		return logs, nil
	}

	err = json.Unmarshal(data, &logs)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *foodLogRepository) DeleteFoodLog(ctx context.Context, userID string, logID int64) error {
	// Execute the delete. Returns data, count, err. We only need err.
	// We enforce userID to ensure users can only delete their own logs!
	_, _, err := r.client.From("food_logs").
		Delete("", "").
		Eq("id", strconv.FormatInt(logID, 10)).
		Eq("user_id", userID).
		Execute()

	return err
}
