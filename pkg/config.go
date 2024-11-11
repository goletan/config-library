package config

import (
	"fmt"

	"github.com/goletan/config/internal/config"
	"go.uber.org/zap"
)

// LoadConfig is a wrapper function to load a configuration.
func LoadConfig(configName string, target interface{}, log *zap.Logger) error {
	// Check if the config is already in cache
	if cachedConfig, found := config.LoadConfigFromCache(configName); found {
		if cachedTarget, ok := cachedConfig.(interface{}); ok {
			*target.(*interface{}) = cachedTarget
			return nil
		}
		return fmt.Errorf("invalid type for cached config")
	}

	// If not in cache, load from file and cache it
	return config.LoadConfig(configName, target, log)
}
