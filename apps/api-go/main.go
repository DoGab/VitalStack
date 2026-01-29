package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/firebase/genkit/go/genkit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// HealthOutput represents the health check response
type HealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok" doc:"Server health status"`
	}
}

// ScanInput represents the scan request body
type ScanInput struct {
	Body struct {
		ImageBase64 string `json:"image_base64,omitempty" doc:"Base64 encoded image data"`
	}
}

// MacroData represents nutritional macro information
type MacroData struct {
	Calories int     `json:"calories" example:"450" doc:"Total calories"`
	Protein  float64 `json:"protein" example:"25.5" doc:"Protein in grams"`
	Carbs    float64 `json:"carbs" example:"45.0" doc:"Carbohydrates in grams"`
	Fat      float64 `json:"fat" example:"15.5" doc:"Fat in grams"`
	Fiber    float64 `json:"fiber" example:"5.0" doc:"Fiber in grams"`
}

// ScanOutput represents the scan response
type ScanOutput struct {
	Body struct {
		FoodName    string    `json:"food_name" example:"Grilled Chicken Salad" doc:"Detected food name"`
		Confidence  float64   `json:"confidence" example:"0.92" doc:"Detection confidence score"`
		Macros      MacroData `json:"macros" doc:"Nutritional macro information"`
		ServingSize string    `json:"serving_size" example:"1 plate (350g)" doc:"Estimated serving size"`
	}
}

func main() {
	// Initialize Genkit
	// This is the base initialization block. Add LLM plugins here later:
	// - Google Gemini: genkit.WithPlugins(googlegenai.NewPlugin(ctx, nil))
	// - OpenAI: genkit.WithPlugins(openai.NewPlugin(ctx, nil))
	ctx := context.Background()
	g := genkit.Init(ctx)
	_ = g // Will be used for AI flows later

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router
	router := gin.Default()

	// Configure CORS for SvelteKit dev server
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Create Huma API with OpenAPI 3.1
	api := humagin.New(router, huma.DefaultConfig("MacroGuard API", "1.0.0"))

	// Register health endpoint
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health Check",
		Description: "Returns the health status of the API server",
		Tags:        []string{"System"},
	}, func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = "ok"
		return resp, nil
	})

	// Register scan placeholder endpoint
	huma.Register(api, huma.Operation{
		OperationID: "scan-placeholder",
		Method:      http.MethodPost,
		Path:        "/scan-placeholder",
		Summary:     "Scan Food Image (Placeholder)",
		Description: "Placeholder endpoint that returns dummy macro data. Will be replaced with actual AI-powered food analysis.",
		Tags:        []string{"Scan"},
	}, func(ctx context.Context, input *ScanInput) (*ScanOutput, error) {
		resp := &ScanOutput{}
		resp.Body.FoodName = "Grilled Chicken Salad"
		resp.Body.Confidence = 0.92
		resp.Body.Macros = MacroData{
			Calories: 450,
			Protein:  35.5,
			Carbs:    25.0,
			Fat:      22.0,
			Fiber:    8.5,
		}
		resp.Body.ServingSize = "1 plate (350g)"
		return resp, nil
	})

	// Get port from environment or default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 MacroGuard API starting on http://localhost:%s", port)
	log.Printf("📚 API Documentation available at http://localhost:%s/docs", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
