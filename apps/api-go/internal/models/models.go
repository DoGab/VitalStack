package models

import "time"

type Profile struct {
	ID        string    `json:"id"`
	Weight    *float64  `json:"weight,omitempty"`
	Height    *float64  `json:"height,omitempty"`
	Age       *int      `json:"age,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type FoodLog struct {
	ID                  int64               `json:"id,omitempty"`
	UserID              *string             `json:"user_id,omitempty"`
	FoodName            string              `json:"food_name"`
	DetectionConfidence float64             `json:"detection_confidence"`
	Calories            int                 `json:"calories"`
	Protein             float64             `json:"protein"`
	Carbs               float64             `json:"carbs"`
	Fat                 float64             `json:"fat"`
	Fiber               float64             `json:"fiber"`
	Ingredients         []FoodLogIngredient `json:"food_log_ingredients,omitempty"`
	CreatedAt           time.Time           `json:"created_at,omitempty"`
}

type FoodLogIngredient struct {
	ID              string    `json:"id,omitempty"`
	FoodLogID       int64     `json:"food_log_id"`
	Name            string    `json:"name"`
	ServingSize     *int      `json:"serving_size,omitempty"`
	ServingQuantity *float64  `json:"serving_quantity,omitempty"`
	ServingUnit     *string   `json:"serving_unit,omitempty"`
	Calories        int       `json:"calories"`
	Protein         float64   `json:"protein"`
	Carbs           float64   `json:"carbs"`
	Fat             float64   `json:"fat"`
	Fiber           float64   `json:"fiber"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}
