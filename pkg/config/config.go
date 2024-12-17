package config

import (
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

// NewConfig creates a new instance of Config with default cache settings.
// It initializes the cache capacity and TTL (time-to-live) with default values.
// Returns a pointer to the newly created Config instance.
func NewConfig() *Config {
	cfg := &Config{}
	cfg.Cache.Capacity = defaultCapacity
	cfg.Cache.TTL = YAMLDuration(defaultTTL)
	return cfg
}

// OverrideFromConfigYAML overrides the provided configuration with values from a YAML file.
// It reads the YAML file specified by the 'file' parameter and unmarshals its content into a temporary Config struct.
// If the YAML file contains non-zero values for certain fields, those values will override the corresponding fields in the provided 'cfg' parameter.
// If the YAML file is not found, it logs a message and returns nil.
//
// Parameters:
//   - cfg: A pointer to the Config struct to be overridden.
//   - file: The path to the YAML file containing the configuration overrides.
//
// Returns:
//   - error: An error if there is an issue reading the file or unmarshalling its content, otherwise nil.
func OverrideFromConfigYAML(cfg *Config, file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		log.Println("config.yaml not found")
		return nil
	}

	var fileCfg Config
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return err
	}

	if fileCfg.Cache.Capacity != 0 {
		cfg.Cache.Capacity = fileCfg.Cache.Capacity
	}
	if fileCfg.Cache.TTL != 0 {
		cfg.Cache.TTL = fileCfg.Cache.TTL
	}
	if fileCfg.Cache.Redis.Addr != "" {
		cfg.Cache.Redis.Addr = fileCfg.Cache.Redis.Addr
	}
	if fileCfg.Cache.Redis.Username != "" {
		cfg.Cache.Redis.Username = fileCfg.Cache.Redis.Username
	}
	if fileCfg.Cache.Redis.Password != "" {
		cfg.Cache.Redis.Password = fileCfg.Cache.Redis.Password
	}
	if fileCfg.Cache.Redis.DB != 0 {
		cfg.Cache.Redis.DB = fileCfg.Cache.Redis.DB
	}

	return nil
}

// OverrideFromEnvironment overrides the configuration values in the provided
// Config struct with values from environment variables, if they are set.
// The following environment variables are checked:
// - CACHE_CAPACITY: sets the Cache.Capacity field (expects an integer value).
// - CACHE_TTL: sets the Cache.TTL field (expects a duration string, e.g., "1h").
// - REDIS_ADDR: sets the Cache.Redis.Addr field (expects a string value).
// - REDIS_USERNAME: sets the Cache.Redis.Username field (expects a string value).
// - REDIS_PASSWORD: sets the Cache.Redis.Password field (expects a string value).
// - REDIS_DB: sets the Cache.Redis.DB field (expects an integer value).
//
// If any of the environment variables contain invalid values, an error is returned.
func OverrideFromEnvironment(cfg *Config) error {
	if v, ok := os.LookupEnv("CACHE_CAPACITY"); ok {
		c, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		cfg.Cache.Capacity = c
	}

	if v, ok := os.LookupEnv("CACHE_TTL"); ok {
		d, err := time.ParseDuration(v)
		if err != nil {
			return err
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
			return err
		}
		cfg.Cache.Redis.DB = db
	}
	return nil
}
