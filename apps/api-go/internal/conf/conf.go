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

	_ = viper.BindPFlags(pflags)
}
