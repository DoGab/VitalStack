package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/dogab/vitalstack/api/internal/conf"
	"github.com/dogab/vitalstack/api/internal/controller"
	"github.com/dogab/vitalstack/api/internal/server"
	"github.com/dogab/vitalstack/api/pkg/service"

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

	// Register nutrition controller (mock or real based on config)
	if viper.GetBool(conf.DevMocksNutritionServiceArg) {
		slog.Info("🧪 Using MOCK nutrition controller")
		api.RegisterAPI(controller.NewNutritionMockController())
	} else {
		// Initialize Genkit
		// This is the base initialization block. Add LLM plugins here later:
		// - Google Gemini: genkit.WithPlugins(googlegenai.NewPlugin(ctx, nil))
		// - OpenAI: genkit.WithPlugins(openai.NewPlugin(ctx, nil))
		g := genkit.Init(serverShutdownContext,
			genkit.WithPlugins(&googlegenai.GoogleAI{}),
			genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
		)
		svc := service.NewNutritionService(g)
		api.RegisterAPI(controller.NewNutritionController(svc))
	}

	// start the server
	err := api.Serve(ctx)
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
