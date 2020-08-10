package lru

import (
	"bytes"
	lru "github.com/hashicorp/golang-lru"
	"io/ioutil"
	"net/http"
	"time"
)

type lruCache struct {
	cache *lru.Cache
}

type item struct {
	uncompressed bool
	statusCode   int
	body         []byte
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

	item, ok := value.(item)
	if !ok {
		return nil, false
	}

	return &http.Response{Uncompressed: item.uncompressed, StatusCode: item.statusCode, Body: ioutil.NopCloser(bytes.NewBuffer(item.body))}, true
}

func (c *lruCache) Set(s string, response *http.Response, duration time.Duration) {
	var bodyBytes []byte
	if response.Body != nil {
		defer response.Body.Close()
		bodyBytes, _ = ioutil.ReadAll(response.Body)
	}
	c.cache.Add(s, item{uncompressed: response.Uncompressed, statusCode: response.StatusCode, body: bodyBytes})
	response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
}
