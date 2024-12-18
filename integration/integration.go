//go:build integration
// +build integration

package integration

import (
	"log"
	"net/http"
	"net/http/httptest"
)

const (
	binaryPath = "./caching-proxy"
)

func mockOriginServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Hello, world!")); err != nil {
			log.Fatalf("failed to write response: %v", err)
		}
	}))
}
