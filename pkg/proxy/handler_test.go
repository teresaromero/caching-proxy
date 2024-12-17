package proxy

import (
	"caching-proxy/pkg/cache"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockCache struct {
	items map[string]*cache.Item
}

func (m *MockCache) Get(_ context.Context, key string) (*cache.Item, bool) {
	item, ok := m.items[key]
	return item, ok
}

func (m *MockCache) Set(key string, item *cache.Item) {
	m.items[key] = item
}

func (m *MockCache) TTL() time.Duration {
	return 1 * time.Minute
}

func TestProxyHandler(t *testing.T) {
	originServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Origin-Header", "ok")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Hello from origin")); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
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

			proxy := &Proxy{
				Origin:     originServer.URL,
				HttpClient: originServer.Client(),
				Cache: &MockCache{
					items: make(map[string]*cache.Item),
				},
			}

			handler := proxy.Handler()
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
