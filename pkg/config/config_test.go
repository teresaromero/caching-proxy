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
`
	if _, err := tempFile.Write([]byte(testData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Override the configFile variable to point to the temp file
	configFile = tempFile.Name()

	// Call the function to test
	cfg, err := ReadFromConfigYAML()
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
}

func TestReadConfigYAMLError(t *testing.T) {
	// Override the configFile variable to point to a non-existent file
	configFile = "non_existent_config.yaml"

	// Call the function to test
	_, err := ReadFromConfigYAML()
	if err == nil {
		t.Fatal("Expected an error when reading a non-existent file, but got nil")
	}
}
func TestReadFromEnvironment(t *testing.T) {
	// Set environment variables for the test
	os.Setenv("CACHE_CAPACITY", "200")
	os.Setenv("CACHE_TTL", "120s")
	defer os.Unsetenv("CACHE_CAPACITY")
	defer os.Unsetenv("CACHE_TTL")

	// Call the function to test
	cfg, err := ReadFromEnvironment()
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
}

func TestReadFromEnvironmentError(t *testing.T) {
	// Set invalid environment variables for the test
	os.Setenv("CACHE_CAPACITY", "invalid")
	os.Setenv("CACHE_TTL", "invalid")
	defer os.Unsetenv("CACHE_CAPACITY")
	defer os.Unsetenv("CACHE_TTL")

	// Call the function to test
	_, err := ReadFromEnvironment()
	if err == nil {
		t.Fatal("Expected an error when reading invalid environment variables, but got nil")
	}
}

