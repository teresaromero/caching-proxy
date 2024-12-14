package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMainFunction(t *testing.T) {
	// Set up a temporary flag set
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"cmd", "-port=8081", "-origin=http://example.com"}

	// Mock the origin server
	originServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Origin-Header", "ok")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from origin"))
	}))
	defer originServer.Close()

	// Override the origin flag to use the mock server
	os.Args = []string{"cmd", "-port=8081", "-origin=" + originServer.URL}

	// Run the main function in a goroutine
	go main()

	// Give the server a moment to start
	// This is not ideal, but sufficient for this simple test
	// In a real-world scenario, you might use synchronization primitives
	// to ensure the server is ready before proceeding
	<-time.After(100 * time.Millisecond)

	// Make a request to the proxy server
	resp, err := http.Get("http://localhost:8081/test")
	if err != nil {
		t.Fatalf("failed to make request to proxy server: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if string(body) != "Hello from origin" {
		t.Errorf("expected body %q, got %q", "Hello from origin", string(body))
	}
	if resp.Header.Get("Origin-Header") != "ok" {
		t.Errorf("expected header %q: %q, got %q: %q", "Origin-Header", "ok", "Origin-Header", resp.Header.Get("Origin-Header"))
	}
	if resp.Header.Get("X-Cache") != "miss" {
		t.Errorf("expected header %q: %q, got %q: %q", "X-Cache", "miss", "X-Cache", resp.Header.Get("X-Cache"))
	}

	// Make a second request to the proxy server
	resp, err = http.Get("http://localhost:8081/test")
	if err != nil {
		t.Fatalf("failed to make request to proxy server: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.Header.Get("X-Cache") != "hit" {
		t.Errorf("expected header %q: %q, got %q: %q", "X-Cache", "miss", "X-Cache", resp.Header.Get("X-Cache"))
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if string(body) != "Hello from origin" {
		t.Errorf("expected body %q, got %q", "Hello from origin", string(body))
	}
	if resp.Header.Get("Origin-Header") != "ok" {
		t.Errorf("expected header %q: %q, got %q: %q", "Origin-Header", "ok", "Origin-Header", resp.Header.Get("Origin-Header"))
	}

}
