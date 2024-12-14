package config

import (
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	defaultCapacity = 100
	defaultTTL      = 5 * time.Minute
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

// readFromConfigYAML reads the configuration from a YAML file specified by the
// configFile variable, unmarshals it into a Config struct, and returns a
// pointer to the Config struct. If there is an error reading the file or
// unmarshalling the YAML, it returns an error.
func readFromConfigYAML() (*Config, error) {
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

// readFromEnvironment reads configuration settings from environment variables
// and returns a Config struct populated with these settings. It looks for the
// following environment variables:
// - CACHE_CAPACITY: an integer representing the cache capacity.
// - CACHE_TTL: a duration string representing the time-to-live for cache entries.
//
// If the environment variables are not set or if there is an error parsing their
// values, it returns an error.
func readFromEnvironment() (*Config, error) {
	cfg := &Config{}
	if v, ok := os.LookupEnv("CACHE_CAPACITY"); ok {
		c, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		cfg.Cache.Capacity = c
	}

	if v, ok := os.LookupEnv("CACHE_TTL"); ok {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, err
		}
		cfg.Cache.TTL = YAMLDuration(d)
	}
	return cfg, nil
}

// Read reads the configuration for the application from multiple sources.
// It initializes the configuration with default values, then overrides them
// with values from a YAML configuration file if available, and finally overrides
// them with values from environment variables if available.
//
// The precedence order for configuration values is:
// 1. Environment variables
// 2. YAML configuration file
// 3. Default values
//
// Returns a pointer to the Config struct and an error if any occurred during reading
// from the YAML configuration file or environment variables.
func Read() (*Config, error) {
	cfg := &Config{}
	cfg.Cache.Capacity = defaultCapacity
	cfg.Cache.TTL = YAMLDuration(defaultTTL)

	cfgYAML, err := readFromConfigYAML()
	if err != nil {
		return nil, err
	}
	if cfgYAML != nil {
		if cfgYAML.Cache.Capacity > 0 {
			cfg.Cache.Capacity = cfgYAML.Cache.Capacity
		}
		if cfgYAML.Cache.TTL > 0 {
			cfg.Cache.TTL = cfgYAML.Cache.TTL
		}
	}

	cfgEnv, err := readFromEnvironment()
	if err != nil {
		return nil, err
	}
	if cfgEnv != nil {
		if cfgEnv.Cache.Capacity > 0 {
			cfg.Cache.Capacity = cfgEnv.Cache.Capacity
		}
		if cfgEnv.Cache.TTL > 0 {
			cfg.Cache.TTL = cfgEnv.Cache.TTL
		}
	}

	return cfg, nil
}
