package proxy

import (
	"caching-proxy/pkg/cache"
	"context"
	"io"
	"log"
	"net/http"
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

		// request to origin server
		originURL := p.Origin + r.URL.Path
		if r.URL.RawQuery != "" {
			originURL += "?" + r.URL.RawQuery
		}

		log.Println("forwarding request to origin server:", originURL)
		req, err := http.NewRequest(r.Method, originURL, io.NopCloser(r.Body))
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
