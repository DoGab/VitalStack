package cmd

import (
	"github.com/dogab/macroguard/api/internal/conf"
	"github.com/dogab/macroguard/api/internal/controller"
	"github.com/dogab/macroguard/api/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var openAPICmd = &cobra.Command{
	Use:   "openapi",
	Short: "Generate OpenAPI specification",
	Long:  "Generate OpenAPI specification for the API and save it to a file",
	RunE:  openapiEntryPoint,
}

func openapiEntryPoint(cmd *cobra.Command, _ []string) error {
	nutritionController := controller.NewNutritionMockController()
	api, _ := server.NewServer(":8080")

	// register API endpoints
	api.RegisterAPI(nutritionController)

	return api.OpenAPI(viper.GetString(conf.OpenAPIPathArg), server.SpecFormat(viper.GetString(conf.OpenAPIFormatArg)))
}
