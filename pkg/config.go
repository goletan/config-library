package config

import (
	"github.com/goletan/config/internal/config"
	"github.com/goletan/observability/shared/logger"
	"go.uber.org/zap"
)

// LoadConfig is a wrapper function to load a configuration.
func LoadConfig[T any](configName string, target *T, log *logger.ZapLogger) error {
	// Check if the config is already in cache
	if cachedConfig, found := config.LoadConfigFromCache[T](configName); found {
		*target = *cachedConfig
		return nil
	}

	// If not in cache, load from file and cache it
	if err := config.LoadConfig(configName, target, log); err != nil {
		log.WithContext(map[string]interface{}{
			"step":    "config loading",
			"error":   zap.Error(err),
			"message": "Failed to load events configuration",
		})
		return err
	}

	// Store the loaded config in cache
	config.StoreConfigInCache(configName, target)

	return nil
}
