// /config/internal/config/cache.go
package config

import "sync"

// ConfigCache to store loaded configurations using sync.Map
var configCache sync.Map

// StoreConfigInCache stores a configuration in the cache.
func StoreConfigInCache(configName string, config interface{}) {
	configCache.Store(configName, config)
}

// LoadConfigFromCache retrieves a configuration from the cache if it exists.
func LoadConfigFromCache[T any](configName string) (*T, bool) {
	if val, ok := configCache.Load(configName); ok {
		cfg, isCorrectType := val.(*T)
		if !isCorrectType {
			return nil, false
		}
		return cfg, true
	}
	return nil, false
}
