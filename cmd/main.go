package main

import (
	cacheable "github.com/TaylorOno/http-cacheable"
	"github.com/TaylorOno/http-cacheable-demo/cmd/routes"
	"github.com/TaylorOno/http-cacheable-demo/internal/lru"
	"github.com/TaylorOno/http-cacheable-demo/internal/lruttl"
	"github.com/TaylorOno/http-cacheable-demo/internal/ttl"
	"net/http"
)

func main() {
	goCacheClient := cacheable.NewCacheableMiddleware(ttl.NewTTLCache(), 5)
	lruCacheClient := cacheable.NewCacheableMiddleware(lru.NewLRUCache(5), 5)
	lruttlCacheClient := cacheable.NewCacheableMiddleware(lruttl.NewLRUTTLCache(5), 5)

	server := routes.Server{
		GoCacheClient:     goCacheClient(&http.Client{}),
		LRUCacheClient:    lruCacheClient(&http.Client{}),
		LRUTTLCacheClient: lruttlCacheClient(&http.Client{}),
	}

	http.HandleFunc("/goCache/", server.GoCache)
	http.HandleFunc("/lruCache/", server.LRUCache)
	http.HandleFunc("/lruttlCache/", server.LRUTTLCache)
	http.ListenAndServe(":8090", nil)
}
