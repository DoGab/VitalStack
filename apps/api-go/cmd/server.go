package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/dogab/macroguard/api/internal/conf"
	"github.com/dogab/macroguard/api/internal/controller"
	"github.com/dogab/macroguard/api/internal/server"
	"github.com/dogab/macroguard/api/pkg/service"

	"github.com/firebase/genkit/go/genkit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerEntryPoint(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	// create context that listens for the interrupt signal from the OS
	serverShutdownContext, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize Genkit
	// This is the base initialization block. Add LLM plugins here later:
	// - Google Gemini: genkit.WithPlugins(googlegenai.NewPlugin(ctx, nil))
	// - OpenAI: genkit.WithPlugins(openai.NewPlugin(ctx, nil))
	g := genkit.Init(ctx)
	_ = g // Will be used for AI flows later

	serverAddr := viper.GetString(conf.ServerAddrArg)

	svc := service.NewNutritionService()
	ctrl := controller.NewNutritionController(svc)

	api, shutdown := server.NewServer(serverAddr)

	// register API endpoints
	api.RegisterAPI(ctrl)

	// start the server
	err := api.Serve(ctx)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	slog.Info("🚀 MacroGuard API starting", "address", serverAddr)
	slog.Info("📚 API Documentation available", "address", serverAddr+"/docs")

	<-serverShutdownContext.Done()

	_, cancel := context.WithTimeout(serverShutdownContext, 5*time.Second)
	defer cancel()

	if err := shutdown(serverShutdownContext); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	slog.Info("MacroGuard API stopped")

	return nil
}
