package main

import (
	cacheable "github.com/TaylorOno/http-cacheable"
	"github.com/TaylorOno/http-cacheable-demo/cmd/routes"
	"github.com/TaylorOno/http-cacheable-demo/internal"
	"net/http"
	"time"
)

func main() {
	ttlCache := internal.NewTTLCache(5*time.Second)
	lruCache := internal.NewLRUCache(5)
	lruttlCache := internal.NewLRUTTLCache(5)
	memCache := internal.NewMemCacheClient()

	goCacheClient := cacheable.NewCacheableMiddleware(ttlCache, 60)
	lruCacheClient := cacheable.NewCacheableMiddleware(lruCache, 60)
	lruttlCacheClient := cacheable.NewCacheableMiddleware(lruttlCache, 60)
	memCacheClient := cacheable.NewCacheableMiddleware(memCache, 60)
	multiStageClient := cacheable.NewCacheableMiddleware(internal.NewMultiStageCache(ttlCache, memCache), 60)


	server := routes.Server{
		GoCacheClient:     goCacheClient(&http.Client{}),
		LRUCacheClient:    lruCacheClient(&http.Client{}),
		LRUTTLCacheClient: lruttlCacheClient(&http.Client{}),
		MemCacheClient:    memCacheClient(&http.Client{}),
		MultiStageClient:  multiStageClient(&http.Client{}),
	}

	http.HandleFunc("/goCache/", server.GoCache)
	http.HandleFunc("/lruCache/", server.LRUCache)
	http.HandleFunc("/lruttlCache/", server.LRUTTLCache)
	http.HandleFunc("/memCached/", server.MemCache)
	http.HandleFunc("/multiStage/", server.MultiStageCache)
	http.ListenAndServe(":8090", nil)
}
