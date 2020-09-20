package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type lruCache struct {
	cache *lru.Cache
}

func NewLRUCache(size int) *lruCache {
	cache, _ := lru.New(size)
	return &lruCache{cache: cache}
}

func (c *lruCache) Get(s string) (*http.Response, bool) {
	value, ok := c.cache.Get(s)
	if !ok {
		return nil, false
	}

	respBytes, ok := value.([]byte)
	if !ok {
		return nil, false
	}

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(respBytes)), nil)
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
		return nil, false
	}

	log.Print("lru get")
	return resp, true
}

func (c *lruCache) Set(s string, response *http.Response, duration time.Duration) {
	responseBytes, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
	}
	log.Print("lru set")
	c.cache.Add(s, responseBytes)
}
