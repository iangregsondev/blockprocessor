package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ggwhite/go-masker/v2"
	"github.com/iangregsondev/deblockprocessor/internal/bitcoinprocessor/config"
	"github.com/iangregsondev/deblockprocessor/internal/bitcoinprocessor/di"
	oswrapper "github.com/iangregsondev/deblockprocessor/internal/wrappers/os"
	loggerwrapper "github.com/iangregsondev/deblockprocessor/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	configPath string
	envPath    string
)

func main() {
	// Create a new slog.Logger with the desired configuration
	slogLogger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				// TODO: Bring over the log level from the config!
				Level: slog.LevelDebug,
			},
		),
	)

	logger := loggerwrapper.NewSlogWrapper(slogLogger)

	// Create new wrapper for OS
	osWrapper := oswrapper.NewOSWrapper()

	mask := masker.NewMaskerMarshaler()

	// Create a new Cobra command
	var rootCmd = &cobra.Command{
		Use:          "myapp",
		Short:        "MyApp is an example application",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Initialize the config
			cfg, err := config.NewConfig(logger, osWrapper, configPath, envPath)
			if err != nil {
				return fmt.Errorf("failed to create new config: %w", err)
			}

			app, err := di.NewDI(logger, mask, osWrapper, cfg)
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
