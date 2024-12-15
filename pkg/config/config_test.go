package config

import (
	"os"
	"testing"
	"time"
)

func TestReadConfigYAML(t *testing.T) {
	// Create a temporary config file
	tempFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data to the temp file
	testData := `
cache:
  capacity: 100
  ttl: 60s
  redis:
    addr: localhost:6379
    username: "admin"
    password: "password"
    db: 0
`
	if _, err := tempFile.Write([]byte(testData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Override the configFile variable to point to the temp file
	configFile := tempFile.Name()

	// Call the function to test
	cfg, err := readFromConfigYAML(configFile)
	if err != nil {
		t.Fatalf("ReadConfigYAML() returned an error: %v", err)
	}

	// Validate the results
	if cfg.Cache.Capacity != 100 {
		t.Errorf("Expected Cache.Capacity to be 100, got %d", cfg.Cache.Capacity)
	}
	if cfg.Cache.TTL != YAMLDuration(time.Duration(60)*time.Second) {
		t.Errorf("Expected Cache.TTL to be 60s, got %v", cfg.Cache.TTL)
	}
	if cfg.Cache.Redis.Addr != "localhost:6379" {
		t.Errorf("Expected Cache.Redis.Addr to be localhost:6379, got %s", cfg.Cache.Redis.Addr)
	}
	if cfg.Cache.Redis.Username != "admin" {
		t.Errorf("Expected Cache.Redis.Username to be admin, got %s", cfg.Cache.Redis.Username)
	}
	if cfg.Cache.Redis.Password != "password" {
		t.Errorf("Expected Cache.Redis.Password to be password, got %s", cfg.Cache.Redis.Password)
	}
	if cfg.Cache.Redis.DB != 0 {
		t.Errorf("Expected Cache.Redis.DB to be 0, got %d", cfg.Cache.Redis.DB)
	}
}

func TestReadFromEnvironment(t *testing.T) {
	// Set environment variables for the test
	t.Setenv("CACHE_CAPACITY", "200")
	t.Setenv("CACHE_TTL", "120s")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("REDIS_USERNAME", "admin")
	t.Setenv("REDIS_PASSWORD", "password")
	t.Setenv("REDIS_DB", "1")

	// Call the function to test
	cfg, err := readFromEnvironment()
	if err != nil {
		t.Fatalf("ReadFromEnvironment() returned an error: %v", err)
	}

	// Validate the results
	if cfg.Cache.Capacity != 200 {
		t.Errorf("Expected Cache.Capacity to be 200, got %d", cfg.Cache.Capacity)
	}
	if cfg.Cache.TTL != YAMLDuration(time.Duration(120)*time.Second) {
		t.Errorf("Expected Cache.TTL to be 120s, got %v", cfg.Cache.TTL)
	}
	if cfg.Cache.Redis.Addr != "localhost:6379" {
		t.Errorf("Expected Cache.Redis.Addr to be localhost:6379, got %s", cfg.Cache.Redis.Addr)
	}
	if cfg.Cache.Redis.Username != "admin" {
		t.Errorf("Expected Cache.Redis.Username to be admin, got %s", cfg.Cache.Redis.Username)
	}
	if cfg.Cache.Redis.Password != "password" {
		t.Errorf("Expected Cache.Redis.Password to be password, got %s", cfg.Cache.Redis.Password)
	}
	if cfg.Cache.Redis.DB != 1 {
		t.Errorf("Expected Cache.Redis.DB to be 1, got %d", cfg.Cache.Redis.DB)
	}
}

func TestReadFromEnvironmentError(t *testing.T) {
	// Set invalid environment variables for the test
	os.Setenv("CACHE_CAPACITY", "invalid")
	os.Setenv("CACHE_TTL", "invalid")
	defer os.Unsetenv("CACHE_CAPACITY")
	defer os.Unsetenv("CACHE_TTL")

	// Call the function to test
	_, err := readFromEnvironment()
	if err == nil {
		t.Fatal("Expected an error when reading invalid environment variables, but got nil")
	}
}
