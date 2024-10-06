// config/config.go
package config

import (
	"errors"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// LoadConfig loads configuration from a file into the provided target struct.
func LoadConfig(configName string, paths []string, target interface{}, log *zap.Logger) error {
	if configName == "" {
		return errors.New("configuration file name is empty")
	}

	logger := log
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")

	// Add each provided path for viper to search for the config file
	for _, path := range paths {
		viper.AddConfigPath(path)
	}

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Failed to read configuration file", zap.Error(err))
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Unmarshal the configuration into the target struct
	if err := viper.Unmarshal(target); err != nil {
		logger.Error("Failed to parse configuration", zap.Error(err))
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Set up hot-reloading
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Info("Configuration file changed", zap.String("file", e.Name))

		if err := viper.Unmarshal(target); err != nil {
			logger.Error("Failed to reload configuration", zap.Error(err))
		} else {
			logger.Info("Configuration reloaded successfully")
		}
	})
	viper.WatchConfig()

	logger.Info("Configuration loaded successfully")
	return nil
}
