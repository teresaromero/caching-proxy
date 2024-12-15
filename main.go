package main

import (
	"caching-proxy/pkg/cache"
	"caching-proxy/pkg/config"
	"caching-proxy/pkg/proxy"
	"context"
	"flag"
	"log"
	"net/http"
	"time"
)

func main() {
	port := flag.String("port", "8080", "port to listen on")
	origin := flag.String("origin", "", "origin host")
	clearCache := flag.Bool("clear-cache", false, "clear cache")
	flag.Parse()

	if *origin == "" && !*clearCache {
		log.Fatal("origin server URL is required")
	}

	config, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	cache := cache.New(
		&cache.CacheConfig{
			TTL:           time.Duration(config.Cache.TTL),
			Capacity:      config.Cache.Capacity,
			RedisAddr:     config.Cache.Redis.Addr,
			RedisDB:       config.Cache.Redis.DB,
			RedisPwd:      config.Cache.Redis.Password,
			RedisUsername: config.Cache.Redis.Username,
		})

	if *clearCache {
		if err := cache.RemoveAll(context.Background()); err != nil {
			log.Fatal(err)
			return
		}
		log.Println("Cache cleared")
		return
	}

	proxy := proxy.Proxy{
		Origin:     *origin,
		HttpClient: &http.Client{},
		Cache:      cache,
	}

	log.Printf("ListenAndServe on port %s ...", *port)
	http.HandleFunc("/", proxy.Handler())
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal(err)
	}
}
