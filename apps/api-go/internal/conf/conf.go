package conf

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// Logging
	loggingKey = "logging."
	// LoggingLevelArg is the flag name for the log level
	LoggingLevelArg = loggingKey + "level"
	// LoggingLevelDefault is the default log level
	LoggingLevelDefault = "info"
	// LoggingLevelHelp is the help message for the log level flag
	LoggingLevelHelp = "Set the log level"

	// LoggingEncodingArg is the flag name for the log encoding
	LoggingEncodingArg = loggingKey + "encoding"
	// LoggingEncodingDefault is the default log encoding
	LoggingEncodingDefault = "json"
	// LoggingEncodingHelp is the help message for the log encoding flag
	LoggingEncodingHelp = "Set the log encoding format"

	// Server
	serverKey = "server."
	// ServerAddrArg is the flag name for the server address
	ServerAddrArg = serverKey + "addr"
	// ServerAddrDefault is the default server address
	ServerAddrDefault = "localhost:8080"
	// ServerAddrHelp is the help message for the server address flag
	ServerAddrHelp = "Set the server address (format: host:port)"

	// ServerOriginArg is the flag name for the server origin
	ServerOriginArg = serverKey + "origin"
	// ServerOriginHelp is the help message for the server origin flag
	ServerOriginHelp = "A list of origins for the server (format: host:port), multiple origins can be set"

	// ServerShutdownTimeoutArg is the flag name for the server shutdown grace period
	ServerShutdownTimeoutArg = serverKey + "shutdown-timeout"
	// ServerShutdownTimeoutDefault is the default server shutdown grace period
	ServerShutdownTimeoutDefault = 60 * time.Second
	// ServerShutdownTimeoutHelp is the help message for the server shutdown grace period
	ServerShutdownTimeoutHelp = "Duration to wait for the server to shutdown gracefully"

	// OpenAPI
	openapiKey = "openapi."
	// OpenAPIPathArg is the flag name for the OpenAPI path
	OpenAPIPathArg = openapiKey + "path"
	// OpenAPIPathDefault is the default OpenAPI path
	OpenAPIPathDefault = "openapi.yaml"
	// OpenAPIPathHelp is the help message for the OpenAPI path flag
	OpenAPIPathHelp = "Set the OpenAPI path"

	// OpenAPIFormatArg is the flag name for the OpenAPI format
	OpenAPIFormatArg = openapiKey + "format"
	// OpenAPIFormatDefault is the default OpenAPI format
	OpenAPIFormatDefault = "yaml"
	// OpenAPIFormatHelp is the help message for the OpenAPI format flag
	OpenAPIFormatHelp = "Set the OpenAPI format"

	// Dev Mocks
	devKey = "dev."
	// DevModeEnabledArg is the flag name for enabling dev mode
	DevModeEnabledArg = devKey + "mode.enabled"
	// DevModeEnabledDefault is the default value for dev mode
	DevModeEnabledDefault = false
	// DevModeEnabledHelp is the help message for the dev mode flag
	DevModeEnabledHelp = "Enable dev mode"

	// DevMocksNutritionServiceArg is the flag name for enabling mock nutrition service
	DevMocksNutritionServiceArg = devKey + "mocks.nutrition-service"
	// DevMocksNutritionServiceDefault is the default value for mock nutrition service
	DevMocksNutritionServiceDefault = false
	// DevMocksNutritionServiceHelp is the help message for the mock nutrition service flag
	DevMocksNutritionServiceHelp = "Enable mock nutrition service for testing"

	// DevMocksScanFoodArg is the flag name for mocking *only* the Genkit scanner
	DevMocksScanFoodArg = devKey + "mocks.scan-food"
	// DevMocksScanFoodDefault is the default value for mock scan food
	DevMocksScanFoodDefault = false
	// DevMocksScanFoodHelp is the help message for the mock scan food flag
	DevMocksScanFoodHelp = "Enable mock AI response for image scanning"

	// Supabase
	supabaseKey = "supabase."
	// SupabaseURLArg is the flag name for the Supabase URL
	SupabaseURLArg = supabaseKey + "url"
	// SupabaseURLDefault is the default value for Supabase URL
	SupabaseURLDefault = ""
	// SupabaseURLHelp is the help message for Supabase URL
	SupabaseURLHelp = "Supabase API URL"

	// SupabaseServiceKeyArg is the flag name for the Supabase Service Key
	SupabaseServiceKeyArg = supabaseKey + "key"
	// SupabaseServiceKeyDefault is the default value for Supabase Service Key
	SupabaseServiceKeyDefault = ""
	// SupabaseServiceKeyHelp is the help message for Supabase Service Key
	SupabaseServiceKeyHelp = "Supabase Service Role Key (for backend server access)"

	// Meilisearch
	meilisearchKey = "meilisearch."
	// MeilisearchURLArg is the flag name for the Meilisearch URL
	MeilisearchURLArg = meilisearchKey + "url"
	// MeilisearchURLDefault is the default value for Meilisearch URL
	MeilisearchURLDefault = "http://localhost:7700"
	// MeilisearchURLHelp is the help message for Meilisearch URL
	MeilisearchURLHelp = "Meilisearch server URL"

	// MeilisearchAPIKeyArg is the flag name for the Meilisearch API key
	MeilisearchAPIKeyArg = meilisearchKey + "api-key"
	// MeilisearchAPIKeyDefault is the default value for Meilisearch API key
	MeilisearchAPIKeyDefault = ""
	// MeilisearchAPIKeyHelp is the help message for Meilisearch API key
	MeilisearchAPIKeyHelp = "Meilisearch API key (derived from master key in production)"

	// USDA FoodData Central
	usdaKey = "usda."
	// USDAAPIKeyArg is the flag name for the USDA API key
	USDAAPIKeyArg = usdaKey + "api-key"
	// USDAAPIKeyDefault is the default value for the USDA API key
	USDAAPIKeyDefault = ""
	// USDAAPIKeyHelp is the help message for the USDA API key
	USDAAPIKeyHelp = "USDA FoodData Central API key (free at https://api.data.gov/signup/)" //nolint:gosec // Not a credential, just a help string

	// Open Food Facts
	offKey = "openfoodfacts."
	// OFFBaseURLArg is the flag name for the OFF Base URL
	OFFBaseURLArg = offKey + "base-url"
	// OFFBaseURLDefault is the default value for the OFF Base URL
	OFFBaseURLDefault = "https://ch.openfoodfacts.org"
	// OFFBaseURLHelp is the help message for the OFF Base URL
	OFFBaseURLHelp = "Open Food Facts server Base URL"

	// OFFLanguageArg is the flag name for the OFF language
	OFFLanguageArg = offKey + "language"
	// OFFLanguageDefault is the default value for the OFF language
	OFFLanguageDefault = "en"
	// OFFLanguageHelp is the help message for the OFF language
	OFFLanguageHelp = "Open Food Facts response language code (e.g. en, fr, de)"

	// OFFSortByArg is the flag name for the OFF sort-by
	OFFSortByArg = offKey + "sort-by"
	// OFFSortByDefault is the default value for the OFF sort-by
	OFFSortByDefault = "popularity"
	// OFFSortByHelp is the help message for the OFF sort-by
	OFFSortByHelp = "Open Food Facts default sorting (e.g. popularity)"
)

var (
	// ServerOriginDefault is the default server origin
	ServerOriginDefault = []string{"http://localhost:3000"}
)

func RegisterFlags(cmd *cobra.Command) {
	pflags := cmd.PersistentFlags()

	// Logging
	pflags.String(LoggingLevelArg, LoggingLevelDefault, LoggingLevelHelp)
	pflags.String(LoggingEncodingArg, LoggingEncodingDefault, LoggingEncodingHelp)

	// Server
	pflags.String(ServerAddrArg, ServerAddrDefault, ServerAddrHelp)
	pflags.StringSlice(ServerOriginArg, ServerOriginDefault, ServerOriginHelp)
	pflags.Duration(ServerShutdownTimeoutArg, ServerShutdownTimeoutDefault, ServerShutdownTimeoutHelp)

	// OpenAPI
	pflags.String(OpenAPIPathArg, OpenAPIPathDefault, OpenAPIPathHelp)
	pflags.String(OpenAPIFormatArg, OpenAPIFormatDefault, OpenAPIFormatHelp)

	// Dev Mocks
	pflags.Bool(DevModeEnabledArg, DevModeEnabledDefault, DevModeEnabledHelp)
	pflags.Bool(DevMocksNutritionServiceArg, DevMocksNutritionServiceDefault, DevMocksNutritionServiceHelp)
	pflags.Bool(DevMocksScanFoodArg, DevMocksScanFoodDefault, DevMocksScanFoodHelp)

	// Supabase
	pflags.String(SupabaseURLArg, SupabaseURLDefault, SupabaseURLHelp)
	pflags.String(SupabaseServiceKeyArg, SupabaseServiceKeyDefault, SupabaseServiceKeyHelp)

	// Meilisearch
	pflags.String(MeilisearchURLArg, MeilisearchURLDefault, MeilisearchURLHelp)
	pflags.String(MeilisearchAPIKeyArg, MeilisearchAPIKeyDefault, MeilisearchAPIKeyHelp)

	// USDA
	pflags.String(USDAAPIKeyArg, USDAAPIKeyDefault, USDAAPIKeyHelp)

	// Open Food Facts
	pflags.String(OFFBaseURLArg, OFFBaseURLDefault, OFFBaseURLHelp)
	pflags.String(OFFLanguageArg, OFFLanguageDefault, OFFLanguageHelp)
	pflags.String(OFFSortByArg, OFFSortByDefault, OFFSortByHelp)

	_ = viper.BindPFlags(pflags)
}
