//go:build integration
// +build integration

package integration

import (
	"net/http"
	"os/exec"
	"testing"
	"time"
)

func TestIntegration_DefaultConfig(t *testing.T) {
	mockOrigin := mockOriginServer()

	// Run the binary
	cmdRun := exec.Command(binaryPath, "--origin", mockOrigin.URL)
	if err := cmdRun.Start(); err != nil {
		t.Fatalf("failed to start caching-proxy: %v", err)
	}
	defer func() {
		if err := cmdRun.Process.Kill(); err != nil {
			t.Fatalf("failed to kill caching-proxy: %v", err)
		}
	}()

	// Wait a bit for the server to start listening
	time.Sleep(2 * time.Second)

	// Test the endpoint
	resp, err := http.Get("http://localhost:8080/some-endpoint")
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if resp.Header.Get("X-Cache") != "miss" {
		t.Fatalf("expected cache miss, got %s", resp.Header.Get("X-Cache"))
	}

	// Test the endpoint again
	resp, err = http.Get("http://localhost:8080/some-endpoint")
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if resp.Header.Get("X-Cache") != "hit" {
		t.Fatalf("expected cache hit, got %s", resp.Header.Get("X-Cache"))
	}

}
