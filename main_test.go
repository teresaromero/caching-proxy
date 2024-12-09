package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestProxyHandler(t *testing.T) {
	originServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Origin-Header", "ok")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from origin"))
	}))
	defer originServer.Close()

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
		expectedBody   string
		expectedHeader http.Header
	}{
		{
			name:           "request",
			method:         http.MethodGet,
			path:           "/test",
			body:           "",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello from origin",
			expectedHeader: http.Header{"Origin-Header": []string{"ok"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler := proxyHandler(originServer.URL, originServer.Client())
			handler.ServeHTTP(w, req)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
			if string(body) != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, string(body))
			}
			for k, v := range tt.expectedHeader {
				if resp.Header.Get(k) != v[0] {
					t.Errorf("expected header %q: %q, got %q: %q", k, v[0], k, resp.Header.Get(k))
				}
			}
		})
	}
}

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
}
