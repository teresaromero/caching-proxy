package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var (
	configFile = "config.yaml"
)

// Config represents the configuration settings for the caching proxy.
// It contains settings related to the cache, including its capacity and TTL (time-to-live).
type Config struct {
	// Cache holds the cache-specific configuration settings.
	Cache struct {
		// Capacity defines the maximum number of items the cache can hold.
		Capacity int `yaml:"capacity"`

		// TTL specifies the duration for which an item should remain in the cache.
		TTL YAMLDuration `yaml:"ttl"`
	} `yaml:"cache"`
}

// ReadConfigYAML reads the configuration from a YAML file specified by the
// configFile variable, unmarshals it into a Config struct, and returns a
// pointer to the Config struct. If there is an error reading the file or
// unmarshalling the YAML, it returns an error.
func ReadConfigYAML() (*Config, error) {
	b, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
