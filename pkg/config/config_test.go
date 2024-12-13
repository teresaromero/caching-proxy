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
	cfg, err := ReadConfigYAML()
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
	_, err := ReadConfigYAML()
	if err == nil {
		t.Fatal("Expected an error when reading a non-existent file, but got nil")
	}
}
