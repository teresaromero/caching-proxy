package proxy

import (
	"io"
	"log"
	"net/http"
)

type Proxy struct {
	Origin     string
	HttpClient *http.Client
}

// Handler returns a http.HandlerFunc that forwards the request to origin server and forwards the response to client
func (p *Proxy) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("incoming new request:", r.Method, r.Host, r.URL.Path)

		// request to origin server
		originURL := p.Origin + r.URL.Path
		if r.URL.RawQuery != "" {
			originURL += "?" + r.URL.RawQuery
		}

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

		// response to client
		for k, v := range originResponse.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(originResponse.StatusCode)
		if originResponse.Body != nil {
			if _, err := io.Copy(w, originResponse.Body); err != nil {
				log.Println("error: parsing origin response body to response", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}
	}
}
