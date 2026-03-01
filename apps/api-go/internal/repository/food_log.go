package repository

import (
	"context"

	"encoding/json"
	"strconv"
	"time"

	"github.com/dogab/vitalstack/api/internal/models"
	"github.com/supabase-community/supabase-go"
)

type FoodLogRepository interface {
	CreateFoodLog(ctx context.Context, log *models.FoodLog) error
	CreateFoodLogIngredient(ctx context.Context, ingredient *models.FoodLogIngredient) error
	GetDailyFoodLogs(ctx context.Context, userID string, startOfDay time.Time, endOfDay time.Time) ([]models.FoodLog, error)
	DeleteFoodLog(ctx context.Context, userID string, logID int64) error
}

type foodLogRepository struct {
	client *supabase.Client
}

func NewFoodLogRepository(client *supabase.Client) FoodLogRepository {
	return &foodLogRepository{client: client}
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
