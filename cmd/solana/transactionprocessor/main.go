package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ggwhite/go-masker/v2"
	"github.com/iangregsondev/deblockprocessor/internal/solana/transactionprocessor/config"
	"github.com/iangregsondev/deblockprocessor/internal/solana/transactionprocessor/di"
	iowrapper "github.com/iangregsondev/deblockprocessor/internal/wrappers/io"
	loggerwrapper "github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
	oswrapper "github.com/iangregsondev/deblockprocessor/internal/wrappers/os"
	"github.com/spf13/cobra"
)

var (
	configPath string
	envPath    string
)

func main() {
	// Define the app name for logger
	appName := "solana-transaction-processor"

	// Create a new slog.Logger
	var logLevel slog.LevelVar

	logLevel.Set(slog.LevelInfo)
	slogLogger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: &logLevel,
				ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						// Add the app name prefix at the start of the log message
						return slog.Attr{Key: "app", Value: slog.StringValue(appName)}
					}

					return a
				},
			},
		),
	)

	logger := loggerwrapper.NewSlogWrapper(slogLogger, &logLevel)

	// Create wrappers
	osWrapper := oswrapper.NewOSWrapper()
	ioWrapper := iowrapper.NewIOWrapper()

	mask := masker.NewMaskerMarshaler()

	// Create a new Cobra command
	var rootCmd = &cobra.Command{
		Use:          "solana-transaction-processor",
		Short:        "Solana transaction processor",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Initialize the config
			cfg, err := config.NewConfig(logger, osWrapper, configPath, envPath)
			if err != nil {
				return fmt.Errorf("failed to create new config: %w", err)
			}

			app, err := di.NewDI(logger, mask, ioWrapper, osWrapper, cfg)
			if err != nil {
				return fmt.Errorf("failed to create new DI: %w", err)
			}

			// Start the application
			if err := app.Run(); err != nil {
				return fmt.Errorf("application error: %w", err)
			}

			return nil
		},
	}

	// Set up the --config flag to specify a configuration file path
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "./", "location of config file (default is ./)")
	// Set up the --env flag to specify an env file path
	rootCmd.PersistentFlags().StringVar(&envPath, "env", "", "location of env file (optional)")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Error executing root command", "error", err)
		os.Exit(1)
	}

	logger.Info("Application terminated successfully.")
	os.Exit(0)
}
