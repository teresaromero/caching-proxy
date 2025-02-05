package proxy

import (
	"caching-proxy/internal/cache"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CacheInterface interface {
	Get(ctx context.Context, key string) (*cache.Item, bool)
	Set(key string, item *cache.Item)
	TTL() time.Duration
}

type Proxy struct {
	Origin     string
	HttpClient *http.Client
	Cache      CacheInterface
}

// Handler returns a http.HandlerFunc that forwards the request to origin server and forwards the response to client
func (p *Proxy) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("incoming new request:", r.Method, r.Host, r.URL.Path)
		ctx := r.Context()
		cacheKey := r.Method + r.Host + r.URL.Path

		// check cache
		if item, ok := p.Cache.Get(ctx, cacheKey); ok {
			for k, v := range item.ResponseHeaders {
				for _, vv := range v {
					w.Header().Add(k, vv)
				}
			}
			w.Header().Add("X-Cache", "hit")
			w.WriteHeader(item.ResponseStatusCode)
			if _, err := w.Write(item.ResponseBody); err != nil {
				log.Println("error: writing cache response to client", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		originURL, err := parseOriginURL(p.Origin)
		if err != nil {
			log.Println("error: parsing origin url", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// request to origin server
		url := originURL + r.URL.Path
		if r.URL.RawQuery != "" {
			url += "?" + r.URL.RawQuery
		}

		log.Println("forwarding request to origin server:", url)
		req, err := http.NewRequest(r.Method, url, io.NopCloser(r.Body))
		req.Header = r.Header.Clone()
		if err != nil {
			log.Println("error: new request forward", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		originResponse, err := p.HttpClient.Do(req)
		if err != nil {
			log.Println("error: request to origin server", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer originResponse.Body.Close()

		body, err := io.ReadAll(originResponse.Body)
		if err != nil {
			log.Println("error: reading origin response body", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// save into cache
		p.Cache.Set(cacheKey, &cache.Item{
			ResponseBody:       body,
			ResponseHeaders:    originResponse.Header,
			ResponseStatusCode: originResponse.StatusCode,
			Expiration:         time.Now().Add(p.Cache.TTL()),
		})

		// response to client
		for k, v := range originResponse.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.Header().Add("X-Cache", "miss")
		w.WriteHeader(originResponse.StatusCode)
		if _, err := w.Write(body); err != nil {
			log.Println("error: writing origin response to client", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func parseOriginURL(origin string) (string, error) {
	// if origin is hostname:port, add default scheme http for url.Parse recognize it as url
	if !strings.Contains(origin, "://") && strings.Contains(origin, ":") {
		origin = "http://" + origin
	}

	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		return "", fmt.Errorf("invalid origin: %s", err)
	}

	if parsedOrigin.Scheme == "" {
		parsedOrigin.Scheme = "http"
	}

	return parsedOrigin.String(), nil
}
