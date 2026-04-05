package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/dogab/vitalstack/api/internal/conf"
	"github.com/dogab/vitalstack/api/internal/controller"
	"github.com/dogab/vitalstack/api/internal/repository"
	"github.com/dogab/vitalstack/api/internal/server"
	"github.com/dogab/vitalstack/api/pkg/datasource"
	"github.com/dogab/vitalstack/api/pkg/search"
	"github.com/dogab/vitalstack/api/pkg/service"
	"github.com/supabase-community/supabase-go"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerEntryPoint(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	// create context that listens for the interrupt signal from the OS
	serverShutdownContext, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverAddr := viper.GetString(conf.ServerAddrArg)

	// Get CORS and dev mode configuration
	allowedOrigins := viper.GetStringSlice(conf.ServerOriginArg)
	devMode := viper.GetBool(conf.DevModeEnabledArg)

	api, shutdown := server.NewServer(
		serverAddr,
		server.WithAllowedOrigins(allowedOrigins),
		server.WithDevMode(devMode),
	)

	supabaseClient, err := supabase.NewClient(viper.GetString(conf.SupabaseURLArg), viper.GetString(conf.SupabaseServiceKeyArg), nil)
	if err != nil {
		slog.Warn("Failed to initialize Supabase client", "error", err)
	}
	foodLogRepo := repository.NewFoodLogRepository(supabaseClient)

	var ctrl server.Controller
	// Register nutrition controller (mock or real based on config)
	if viper.GetBool(conf.DevModeEnabledArg) && viper.GetBool(conf.DevMocksNutritionServiceArg) {
		slog.Info("🧪 Using MOCK nutrition controller")
		ctrl = controller.NewNutritionMockController()
	} else {
		// Initialize Genkit
		g := genkit.Init(serverShutdownContext,
			genkit.WithPlugins(&googlegenai.GoogleAI{}),
			genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
		)

		mockScan := viper.GetBool(conf.DevMocksScanFoodArg)
		svc := service.NewNutritionService(g, foodLogRepo, service.WithMockScan(mockScan))
		ctrl = controller.NewNutritionController(svc)
	}

	// --- Product search & barcode wiring ---
	var productCtrl *controller.ProductController

	meiliClient, meiliErr := search.NewMeilisearchClient(
		viper.GetString(conf.MeilisearchURLArg),
		viper.GetString(conf.MeilisearchAPIKeyArg),
	)
	if meiliErr != nil {
		slog.Warn("Product search unavailable: Meilisearch init failed", "error", meiliErr)
	} else {
		offClient := datasource.NewOFFClient(
			http.DefaultClient,
			datasource.WithBaseURL(viper.GetString(conf.OFFBaseURLArg)),
			datasource.WithLanguage(viper.GetString(conf.OFFLanguageArg)),
			datasource.WithSortBy(viper.GetString(conf.OFFSortByArg)),
		)
		fsvoClient := datasource.NewFSVOClient(
			http.DefaultClient,
			viper.GetString(conf.FSVOBaseURLArg),
			datasource.WithFSVOLanguage(viper.GetString(conf.FSVOLanguageArg)),
		)
		usdaClient := datasource.NewUSDAClient(http.DefaultClient, viper.GetString(conf.USDAAPIKeyArg))
		productSvc := service.NewProductService(meiliClient, offClient, fsvoClient, usdaClient)
		productCtrl = controller.NewProductController(productSvc)
		slog.Info("🔍 Product search enabled")
	}

	// register the endpoints of the controllers
	if productCtrl != nil {
		api.RegisterAPI(ctrl, productCtrl)
	} else {
		api.RegisterAPI(ctrl)
	}

	// start the server
	err = api.Serve(ctx)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	slog.Info("🚀 VitalStack API starting", "address", serverAddr)
	slog.Info("📚 API Documentation available", "address", serverAddr+"/docs")

	<-serverShutdownContext.Done()

	_, cancel := context.WithTimeout(serverShutdownContext, 5*time.Second)
	defer cancel()

	if err := shutdown(serverShutdownContext); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	slog.Info("VitalStack API stopped")

	return nil
}
