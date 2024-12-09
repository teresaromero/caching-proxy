package main

import (
	"flag"
	"io"
	"log"
	"net/http"
)

// proxyHandler returns a http.HandlerFunc that forwards the request to origin server and forwards the response to client
func proxyHandler(origin string, httpClient *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("incoming new request:", r.Method, r.Host, r.URL.Path)

		// request to origin server
		originURL := origin + r.URL.Path
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
		originResponse, err := httpClient.Do(req)
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

func main() {
	port := flag.String("port", "8080", "port to listen on")
	origin := flag.String("origin", "", "origin host")
	flag.Parse()

	if *origin == "" {
		log.Fatal("origin server URL is required")
	}

	log.Printf("ListenAndServe on port %s ...", *port)
	http.HandleFunc("/", proxyHandler(*origin, &http.Client{}))
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal(err)
	}
}
