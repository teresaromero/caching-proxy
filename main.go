package main

import (
	"caching-proxy/pkg/cache"
	"caching-proxy/pkg/config"
	"caching-proxy/pkg/proxy"
	"flag"
	"log"
	"net/http"
	"time"
)

func main() {
	port := flag.String("port", "8080", "port to listen on")
	origin := flag.String("origin", "", "origin host")
	flag.Parse()

	if *origin == "" {
		log.Fatal("origin server URL is required")
	}

	config, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	proxy := proxy.Proxy{
		Origin:     *origin,
		HttpClient: &http.Client{},
		Cache:      cache.New(time.Duration(config.Cache.TTL), config.Cache.Capacity),
	}

	log.Printf("ListenAndServe on port %s ...", *port)
	http.HandleFunc("/", proxy.Handler())
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal(err)
	}
}
