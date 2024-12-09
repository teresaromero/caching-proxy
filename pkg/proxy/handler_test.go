package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

			proxy := &Proxy{
				Origin:     originServer.URL,
				HttpClient: originServer.Client(),
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
