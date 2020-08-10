package lruttl

import (
	"bytes"
	lru "github.com/hashicorp/golang-lru"
	"io/ioutil"
	"net/http"
	"time"
)

type lruttlCache struct {
	cache *lru.Cache
}

type item struct {
	expiration   int64
	uncompressed bool
	statusCode   int
	body         []byte
}

func NewLRUTTLCache(size int) *lruttlCache {
	cache, _ := lru.New(size)
	return &lruttlCache{cache: cache}
}

func (c *lruttlCache) Get(s string) (*http.Response, bool) {
	value, ok := c.cache.Get(s)
	if !ok {
		return nil, false
	}

	item, ok := value.(item)
	if !ok {
		return nil, false
	}

	if time.Now().UnixNano() > item.expiration {
		return nil, false
	}

	return &http.Response{Uncompressed: item.uncompressed, StatusCode: item.statusCode, Body: ioutil.NopCloser(bytes.NewBuffer(item.body))}, true
}

func (c *lruttlCache) Set(s string, response *http.Response, duration time.Duration) {
	var bodyBytes []byte
	if response.Body != nil {
		defer response.Body.Close()
		bodyBytes, _ = ioutil.ReadAll(response.Body)
	}
	c.cache.Add(s, item{expiration: time.Now().Add(duration).UnixNano(), uncompressed: response.Uncompressed, statusCode: response.StatusCode, body: bodyBytes})
	response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
}
