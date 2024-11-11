package config

import "sync"

// ConfigCache to store loaded configurations using sync.Map
var configCache sync.Map

// StoreConfigInCache stores a configuration in the cache.
func StoreConfigInCache(configName string, config interface{}) {
	configCache.Store(configName, config)
}

// LoadConfigFromCache retrieves a configuration from the cache if it exists.
func LoadConfigFromCache(configName string) (interface{}, bool) {
	return configCache.Load(configName)
}
