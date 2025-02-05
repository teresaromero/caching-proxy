package main

import (
	"caching-proxy/internal/cache"
	"caching-proxy/internal/config"
	"caching-proxy/internal/proxy"
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
	configFile := flag.String("config", "config.yaml", "config file")
	flag.Parse()

	if *origin == "" && !*clearCache {
		log.Fatal("origin server URL is required")
	}

	cfg := config.NewConfig()
	if err := config.OverrideFromConfigYAML(cfg, *configFile); err != nil {
		log.Fatal(err)
		return
	}
	if err := config.OverrideFromEnvironment(cfg); err != nil {
		log.Fatal(err)
		return
	}

	cacheInstance := cache.New(
		&cache.CacheConfig{
			TTL:           time.Duration(cfg.Cache.TTL),
			Capacity:      cfg.Cache.Capacity,
			RedisAddr:     cfg.Cache.Redis.Addr,
			RedisDB:       cfg.Cache.Redis.DB,
			RedisPwd:      cfg.Cache.Redis.Password,
			RedisUsername: cfg.Cache.Redis.Username,
		})

	if *clearCache {
		if err := cacheInstance.RemoveAll(context.Background()); err != nil {
			log.Fatal(err)
			return
		}
		log.Println("Cache cleared")
		return
	}

	proxy := proxy.Proxy{
		Origin:     *origin,
		HttpClient: &http.Client{},
		Cache:      cacheInstance,
	}

	log.Printf("ListenAndServe on port %s ...", *port)
	http.HandleFunc("/", proxy.Handler())
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal(err)
	}
}
