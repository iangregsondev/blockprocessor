package config

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iangregsondev/deblockprocessor/internal/solana/blockprocessor/models/config"
	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
	oswrapper "github.com/iangregsondev/deblockprocessor/internal/wrappers/os"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	logger logger.Logger
}

func NewConfig(logger logger.Logger, osWrapper oswrapper.OS, configPath string, envPath string) (*Config, error) {
	// This is the only env file that I will support for now
	envFilename := ".env.local"

	var envFileWithFullPath string

	// Set the path to the .env file
	envFileWithFullPath = path.Join(configPath, envFilename)
	// Load .env file if it exists
	if exists := osWrapper.FileExists(envFileWithFullPath); exists {
		err := godotenv.Load(envFileWithFullPath)
		if err != nil {
			logger.Error("Error loading an .env file", "envFileWithFullPath", envFileWithFullPath, "error", err)
		}
	}

	if envPath != "" {
		// Set the path to the .env file
		envFileWithFullPath = path.Join(envPath, envFilename)
		// Load .env file if it exists
		if exists := osWrapper.FileExists(envFileWithFullPath); exists {
			err := godotenv.Load(path.Join(envPath, envFilename))
			if err != nil {
				logger.Error("Error loading an .env file", "envFileWithFullPath", envFileWithFullPath, "error", err)
			}
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if configPath != "" {
		viper.AddConfigPath(configPath)
	} else {
		viper.AddConfigPath(".")
	}

	appPrefix := "ETHEREUM_BLOCK_PROCESSOR"

	_ = appPrefix
	viper.SetEnvPrefix(appPrefix)
	viper.AutomaticEnv()

	// Configure the key replacer for nested keys
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return &Config{
		logger: logger,
	}, nil
}

func (c *Config) LoadConfig() (*config.Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		c.logger.Error("Error reading config file", "error", err)
	}

	// Unmarshal the configuration into a Config struct
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal configuration: %w", err)
	}

	validate := validator.New()

	err := validate.Struct(cfg)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return nil, fmt.Errorf("invalid validation error: %w", err)
		}

		var validationErrors validator.ValidationErrors
		if !errors.As(err, &validationErrors) {
			return nil, fmt.Errorf("unexecpted validation error: %w", err)
		}

		var errorMessages []string

		for _, err := range validationErrors {
			// Create a custom error message for each validation error
			errorMessage := fmt.Sprintf(
				"Field '%s' (config path: '%s') failed validation: %s",
				err.Field(), err.Namespace(), err.ActualTag(),
			)

			if err.Param() != "" {
				errorMessage += fmt.Sprintf(" (required: %s)", err.Param())
			}

			errorMessages = append(errorMessages, errorMessage)
		}

		// Join the error messages into a single formatted string
		errorMessage := strings.Join(errorMessages, "; ")

		return nil, fmt.Errorf("validation failed: %s", errorMessage)
	}

	return &cfg, nil
}
