package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ConfigCache to store loaded configurations using sync.Map
var configCache sync.Map

// LoadConfig loads configuration from a file into the provided target struct, using sync.Map for efficiency.
func LoadConfig(configName string, target interface{}, log *zap.Logger) error {
	// Check if the config is already in cache
	if cachedConfig, found := configCache.Load(configName); found {
		// Use reflection or type assertion to properly set the value of target
		switch t := target.(type) {
		case *interface{}:
			*t = cachedConfig
		default:
			return fmt.Errorf("unsupported type for target")
		}
		return nil
	}

	viper.SetConfigName(strings.ToLower(configName))
	viper.SetConfigType("yaml")

	// Add common configuration paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Load different configuration files based on environment, with precedence
	loadConfigFiles([]string{
		"./config/override.yaml",
		"./config/tests.yaml",
	})

	envConfigs := []string{"GOLETAN_PROD_CONFIG", "GOLETAN_STAGE_CONFIG", "GOLETAN_LOCAL_CONFIG"}
	for _, envVar := range envConfigs {
		envValue := os.Getenv(envVar)
		if envValue != "" {
			loadConfigFiles([]string{
				fmt.Sprintf("./config/%s.yaml", envValue),
			})
		}
	}

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		if log != nil {
			log.Error("Failed to read configuration file", zap.Error(err))
		}
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Unmarshal the configuration into the target struct
	if err := viper.Unmarshal(target); err != nil {
		if log != nil {
			log.Error("Failed to parse configuration", zap.Error(err))
		}
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Cache the loaded configuration
	configCache.Store(configName, target)

	// Set up hot-reloading
	viper.OnConfigChange(func(e fsnotify.Event) {
		if log != nil {
			log.Info("Configuration file changed", zap.String("file", e.Name))
		}

		if err := viper.Unmarshal(target); err != nil {
			if log != nil {
				log.Error("Failed to reload configuration", zap.Error(err))
			}
		} else {
			// Update the cache after reloading
			configCache.Store(configName, target)

			if log != nil {
				log.Info("Configuration reloaded successfully")
			}
		}
	})
	viper.WatchConfig()

	if log != nil {
		log.Info("Configuration loaded successfully")
	}

	return nil
}

// loadConfigFiles attempts to load a list of configuration files in order, with precedence.
func loadConfigFiles(configFiles []string) {
	for _, configPath := range configFiles {
		if _, err := os.Stat(configPath); err == nil {
			viper.SetConfigFile(configPath)
			if err := viper.MergeInConfig(); err != nil {
				fmt.Printf("Warning: Failed to merge configuration file %s: %v\n", configPath, err)
			}
		}
	}
}
