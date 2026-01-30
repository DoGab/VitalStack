package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/dogab/macroguard/api/internal/conf"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

const (
	configFileArg = "config"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() (err error) {
	ctx := context.Background()

	rootCmd := &cobra.Command{
		Use:   "MacroGuard",
		Short: "MacroGuard API",
		Long:  `MacroGuard is a fullstack sample web application in go`,
		RunE:  ServerEntryPoint,
	}

	// Define the --config flag
	pflags := rootCmd.PersistentFlags()
	pflags.String(configFileArg, "", "config file path")
	_ = viper.BindPFlag(configFileArg, pflags.Lookup(configFileArg))

	// Load configuration after parsing flags
	cobra.OnInitialize(func() {
		if err := initConfig(viper.GetString(configFileArg)); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Register logging
		initLogger()
	})

	// Load configuration flags
	conf.RegisterFlags(rootCmd)

	// Generate markdown documentation
	if len(os.Args) > 1 && os.Args[1] == "gendoc" {
		return doc.GenMarkdownTree(rootCmd, ".")
	}

	return rootCmd.ExecuteContext(ctx)
}

// initconf reads in config file and ENV variables if set.
func initConfig(path string) error {
	if path != "" {
		// Use config file from the flag.
		viper.SetConfigFile(path)

		err := viper.ReadInConfig()
		if err != nil {
			return fmt.Errorf("failed to read configuration file: %w", err)
		}
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	return nil
}

func initLogger() {
	var logger *slog.Logger
	opts := &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: replaceAttr,
	}

	setLogLevel(viper.GetString(conf.LoggingLevelArg), opts)

	switch logEncoding := viper.GetString(conf.LoggingEncodingArg); logEncoding {
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	case "logfmt":
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	default:
		slog.Info("unsupported log encoding, using logfmt", "encoding", logEncoding)
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
	slog.SetDefault(logger)
}

func setLogLevel(level string, opts *slog.HandlerOptions) {
	switch level {
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
	}
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		source, _ := a.Value.Any().(*slog.Source)
		if source != nil {
			source.File = filepath.Base(source.File)
		}
	}
	return a
}
