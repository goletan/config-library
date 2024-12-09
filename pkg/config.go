package config

import (
	"github.com/goletan/config/internal/config"
	"go.uber.org/zap"
)

// LoadConfig is a wrapper function to load a configuration.
func LoadConfig[T any](configName string, target *T, log *zap.Logger) error {
	// Check if the config is already in cache
	if cachedConfig, found := config.LoadConfigFromCache[T](configName); found {
		*target = *cachedConfig
		return nil
	}

	// If not in cache, load from file and cache it
	if err := config.LoadConfig(configName, target, log); err != nil {
		log.Error(
			"Failed to load configuration",
			zap.Error(err),
			zap.Any("context", map[string]interface{}{"step": "config loading"}),
		)
		return err
	}

	// Store the loaded config in cache
	config.StoreConfigInCache(configName, target)

	return nil
}
