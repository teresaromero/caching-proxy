package config

import (
	"errors"
	"log"
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

type Redis struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// Config represents the configuration settings for the caching proxy.
// It contains settings related to the cache, including its capacity and TTL (time-to-live).
type Config struct {
	// Cache holds the cache-specific configuration settings.
	Cache struct {
		// Capacity defines the maximum number of items the cache can hold.
		Capacity int `yaml:"capacity"`

		// TTL specifies the duration for which an item should remain in the cache.
		TTL YAMLDuration `yaml:"ttl"`

		// Redis holds the Redis-specific configuration settings.
		Redis Redis `yaml:"redis"`
	} `yaml:"cache"`
}

// readFromConfigYAML reads the configuration from a YAML file specified by the
// configFile variable, unmarshals it into a Config struct, and returns a
// pointer to the Config struct. If there is an error reading the file or
// unmarshalling the YAML, it returns an error.
func readFromConfigYAML() (*Config, error) {
	b, err := os.ReadFile(configFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("config.yaml not found")
			return nil, nil
		}
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// readFromEnvironment reads configuration values from environment variables
// and populates a Config struct with these values. The following environment
// variables are used:
// - CACHE_CAPACITY: the capacity of the cache (integer).
// - CACHE_TTL: the time-to-live duration for cache entries (duration string).
// - REDIS_ADDR: the address of the Redis server.
// - REDIS_USERNAME: the username for Redis authentication.
// - REDIS_PASSWORD: the password for Redis authentication.
// - REDIS_DB: the Redis database number (integer).
//
// Returns a pointer to a Config struct populated with the values from the
// environment variables, or an error if any of the values cannot be parsed.
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

	if v, ok := os.LookupEnv("REDIS_ADDR"); ok {
		cfg.Cache.Redis.Addr = v
	}
	if v, ok := os.LookupEnv("REDIS_USERNAME"); ok {
		cfg.Cache.Redis.Username = v
	}
	if v, ok := os.LookupEnv("REDIS_PASSWORD"); ok {
		cfg.Cache.Redis.Password = v
	}
	if v, ok := os.LookupEnv("REDIS_DB"); ok {
		db, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		cfg.Cache.Redis.DB = db
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
		if cfgYAML.Cache.Redis.Addr != "" {
			cfg.Cache.Redis.Addr = cfgYAML.Cache.Redis.Addr
		}
		if cfgYAML.Cache.Redis.Username != "" {
			cfg.Cache.Redis.Username = cfgYAML.Cache.Redis.Username
		}
		if cfgYAML.Cache.Redis.Password != "" {
			cfg.Cache.Redis.Password = cfgYAML.Cache.Redis.Password
		}
		if cfgYAML.Cache.Redis.DB > 0 {
			cfg.Cache.Redis.DB = cfgYAML.Cache.Redis.DB
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
		if cfgEnv.Cache.Redis.Addr != "" {
			cfg.Cache.Redis.Addr = cfgEnv.Cache.Redis.Addr
		}
		if cfgEnv.Cache.Redis.Username != "" {
			cfg.Cache.Redis.Username = cfgEnv.Cache.Redis.Username
		}
		if cfgEnv.Cache.Redis.Password != "" {
			cfg.Cache.Redis.Password = cfgEnv.Cache.Redis.Password
		}
		if cfgEnv.Cache.Redis.DB > 0 {
			cfg.Cache.Redis.DB = cfgEnv.Cache.Redis.DB
		}
	}

	return cfg, nil
}
