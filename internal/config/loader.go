package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	logger "github.com/goletan/logger/pkg"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// LoadConfig loads configuration from files into the provided target struct.
func LoadConfig(configName string, target interface{}, log *logger.ZapLogger) error {
	v := viper.New()
	v.SetConfigName(strings.ToLower(configName))
	v.SetConfigType("yaml")

	// Define a more specific configuration path
	// Configurations should now be in the root config directory of each "executable"
	v.AddConfigPath("./config") // Path relative to the executable's root

	// Load environment-specific configuration files with precedence
	envConfigs := []string{"GOLETAN_PROD_CONFIG", "GOLETAN_STAGE_CONFIG", "GOLETAN_LOCAL_CONFIG"}
	for _, envVar := range envConfigs {
		envValue := os.Getenv(envVar)
		if envValue != "" {
			configPath := fmt.Sprintf("./config/%s.yaml", envValue)
			loadConfigFiles([]string{configPath}, v, log)
		}
	}

	// Read the configuration file
	if err := v.ReadInConfig(); err != nil {
		log.Error("Failed to read configuration file", zap.Error(err))
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Unmarshal the configuration into the target struct
	if err := v.Unmarshal(target); err != nil {
		log.Error("Failed to parse configuration", zap.Error(err))
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Set up hot-reloading
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Configuration file changed", zap.String("file", e.Name))

		if err := v.Unmarshal(target); err != nil {
			log.Error("Failed to reload configuration", zap.Error(err))
		} else {
			log.Info("Configuration reloaded successfully")
			// Update the cache after reloading
			StoreConfigInCache(configName, target)
		}
	})
	v.WatchConfig()

	log.Info("Configuration loaded successfully", zap.String("configName", configName))

	// Cache the loaded configuration
	StoreConfigInCache(configName, target)

	return nil
}

// loadConfigFiles attempts to load a list of configuration files in order, with precedence.
func loadConfigFiles(configFiles []string, v *viper.Viper, log *logger.ZapLogger) {
	for _, configPath := range configFiles {
		if _, err := os.Stat(configPath); err == nil {
			v.SetConfigFile(configPath)
			if err := v.MergeInConfig(); err != nil {
				log.Warn("Failed to merge configuration file", zap.String("file", configPath), zap.Error(err))
			}
		}
	}
}
