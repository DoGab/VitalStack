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

	// ServerShutdownTimeoutArg is the flag name for the server shutdown grace period
	ServerShutdownTimeoutArg = serverKey + "shutdown-timeout"
	// ServerShutdownTimeoutDefault is the default server shutdown grace period
	ServerShutdownTimeoutDefault = 60 * time.Second
)

func RegisterFlags(cmd *cobra.Command) {
	pflags := cmd.PersistentFlags()

	// Logging
	pflags.String(LoggingLevelArg, LoggingLevelDefault, LoggingLevelHelp)
	pflags.String(LoggingEncodingArg, LoggingEncodingDefault, LoggingEncodingHelp)

	// Server
	pflags.String(ServerAddrArg, ServerAddrDefault, ServerAddrHelp)
	pflags.Duration(ServerShutdownTimeoutArg, ServerShutdownTimeoutDefault, "Duration to wait for the server to shutdown gracefully")

	_ = viper.BindPFlags(pflags)
}
