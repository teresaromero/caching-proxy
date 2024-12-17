package config

import (
	"os"
	"testing"
	"time"
)

func TestOverrideFromConfigYAML(t *testing.T) {
	tests := []struct {
		name     string
		fileData string
		expected Config
	}{
		{
			name: "Override all fields",
			fileData: `
cache:
  capacity: 200
  ttl: 10m
  redis:
    addr: "localhost:6379"
    username: "user"
    password: "pass"
    db: 1
`,
			expected: Config{
				Cache: struct {
					Capacity int          `yaml:"capacity"`
					TTL      YAMLDuration `yaml:"ttl"`
					Redis    Redis        `yaml:"redis"`
				}{
					Capacity: 200,
					TTL:      YAMLDuration(10 * time.Minute),
					Redis: Redis{
						Addr:     "localhost:6379",
						Username: "user",
						Password: "pass",
						DB:       1,
					},
				},
			},
		},
		{
			name: "Override partial fields",
			fileData: `
cache:
  capacity: 150
  redis:
    addr: "localhost:6380"
`,
			expected: Config{
				Cache: struct {
					Capacity int          `yaml:"capacity"`
					TTL      YAMLDuration `yaml:"ttl"`
					Redis    Redis        `yaml:"redis"`
				}{
					Capacity: 150,
					TTL:      YAMLDuration(0),
					Redis: Redis{
						Addr:     "localhost:6380",
						Username: "",
						Password: "",
						DB:       0,
					},
				},
			},
		},
		{
			name: "No override",
			fileData: `
cache:
  capacity: 0
  ttl: 0
  redis:
    addr: ""
    username: ""
    password: ""
    db: 0
`,
			expected: Config{
				Cache: struct {
					Capacity int          `yaml:"capacity"`
					TTL      YAMLDuration `yaml:"ttl"`
					Redis    Redis        `yaml:"redis"`
				}{
					Capacity: 0,
					TTL:      YAMLDuration(0),
					Redis: Redis{
						Addr:     "",
						Username: "",
						Password: "",
						DB:       0,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "config-*.yaml")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(file.Name())

			if _, err := file.Write([]byte(tt.fileData)); err != nil {
				t.Fatalf("failed to write to temp file: %v", err)
			}
			if err := file.Close(); err != nil {
				t.Fatalf("failed to close temp file: %v", err)
			}

			var cfg Config
			if err := OverrideFromConfigYAML(&cfg, file.Name()); err != nil {
				t.Fatalf("OverrideFromConfigYAML() error = %v", err)
			}

			if cfg.Cache.Capacity != tt.expected.Cache.Capacity {
				t.Errorf("expected capacity %d, got %d", tt.expected.Cache.Capacity, cfg.Cache.Capacity)
			}
			if cfg.Cache.TTL != tt.expected.Cache.TTL {
				t.Errorf("expected ttl %v, got %v", tt.expected.Cache.TTL, cfg.Cache.TTL)
			}
			if cfg.Cache.Redis.Addr != tt.expected.Cache.Redis.Addr {
				t.Errorf("expected redis addr %s, got %s", tt.expected.Cache.Redis.Addr, cfg.Cache.Redis.Addr)
			}
			if cfg.Cache.Redis.Username != tt.expected.Cache.Redis.Username {
				t.Errorf("expected redis username %s, got %s", tt.expected.Cache.Redis.Username, cfg.Cache.Redis.Username)
			}
			if cfg.Cache.Redis.Password != tt.expected.Cache.Redis.Password {
				t.Errorf("expected redis password %s, got %s", tt.expected.Cache.Redis.Password, cfg.Cache.Redis.Password)
			}
			if cfg.Cache.Redis.DB != tt.expected.Cache.Redis.DB {
				t.Errorf("expected redis db %d, got %d", tt.expected.Cache.Redis.DB, cfg.Cache.Redis.DB)
			}
		})
	}
}
func TestOverrideFromEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected Config
	}{
		{
			name: "Override all fields",
			envVars: map[string]string{
				"CACHE_CAPACITY": "200",
				"CACHE_TTL":      "10m",
				"REDIS_ADDR":     "localhost:6379",
				"REDIS_USERNAME": "user",
				"REDIS_PASSWORD": "pass",
				"REDIS_DB":       "1",
			},
			expected: Config{
				Cache: struct {
					Capacity int          `yaml:"capacity"`
					TTL      YAMLDuration `yaml:"ttl"`
					Redis    Redis        `yaml:"redis"`
				}{
					Capacity: 200,
					TTL:      YAMLDuration(10 * time.Minute),
					Redis: Redis{
						Addr:     "localhost:6379",
						Username: "user",
						Password: "pass",
						DB:       1,
					},
				},
			},
		},
		{
			name: "Override partial fields",
			envVars: map[string]string{
				"CACHE_CAPACITY": "150",
				"REDIS_ADDR":     "localhost:6380",
			},
			expected: Config{
				Cache: struct {
					Capacity int          `yaml:"capacity"`
					TTL      YAMLDuration `yaml:"ttl"`
					Redis    Redis        `yaml:"redis"`
				}{
					Capacity: 150,
					TTL:      YAMLDuration(0),
					Redis: Redis{
						Addr:     "localhost:6380",
						Username: "",
						Password: "",
						DB:       0,
					},
				},
			},
		},
		{
			name: "No override",
			envVars: map[string]string{
				"CACHE_CAPACITY": "0",
				"CACHE_TTL":      "0",
				"REDIS_ADDR":     "",
				"REDIS_USERNAME": "",
				"REDIS_PASSWORD": "",
				"REDIS_DB":       "0",
			},
			expected: Config{
				Cache: struct {
					Capacity int          `yaml:"capacity"`
					TTL      YAMLDuration `yaml:"ttl"`
					Redis    Redis        `yaml:"redis"`
				}{
					Capacity: 0,
					TTL:      YAMLDuration(0),
					Redis: Redis{
						Addr:     "",
						Username: "",
						Password: "",
						DB:       0,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			var cfg Config
			if err := OverrideFromEnvironment(&cfg); err != nil {
				t.Fatalf("OverrideReadFromEnvironment() error = %v", err)
			}

			if cfg.Cache.Capacity != tt.expected.Cache.Capacity {
				t.Errorf("expected capacity %d, got %d", tt.expected.Cache.Capacity, cfg.Cache.Capacity)
			}
			if cfg.Cache.TTL != tt.expected.Cache.TTL {
				t.Errorf("expected ttl %v, got %v", tt.expected.Cache.TTL, cfg.Cache.TTL)
			}
			if cfg.Cache.Redis.Addr != tt.expected.Cache.Redis.Addr {
				t.Errorf("expected redis addr %s, got %s", tt.expected.Cache.Redis.Addr, cfg.Cache.Redis.Addr)
			}
			if cfg.Cache.Redis.Username != tt.expected.Cache.Redis.Username {
				t.Errorf("expected redis username %s, got %s", tt.expected.Cache.Redis.Username, cfg.Cache.Redis.Username)
			}
			if cfg.Cache.Redis.Password != tt.expected.Cache.Redis.Password {
				t.Errorf("expected redis password %s, got %s", tt.expected.Cache.Redis.Password, cfg.Cache.Redis.Password)
			}
			if cfg.Cache.Redis.DB != tt.expected.Cache.Redis.DB {
				t.Errorf("expected redis db %d, got %d", tt.expected.Cache.Redis.DB, cfg.Cache.Redis.DB)
			}

			for k := range tt.envVars {
				os.Unsetenv(k)
			}
		})
	}
}
