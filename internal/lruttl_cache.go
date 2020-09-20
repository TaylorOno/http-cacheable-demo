package internal

import (
	"bufio"
	"bytes"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type lruttlCache struct {
	cache *lru.Cache
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

	respBytes, ok := value.([]byte)
	if !ok {
		return nil, false
	}

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(respBytes)), nil)
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
		return nil, false
	}

	expires, err := http.ParseTime(resp.Header.Get("Expires"))
	if err != nil || time.Now().UTC().After(expires) {
		log.Print(fmt.Sprintf("error:%s", "expired soft ttl"))
		return resp, false
	}

	log.Print("lruttl get")
	return resp, true
}

func (c *lruttlCache) Set(s string, response *http.Response, duration time.Duration) {
	if duration > 0 {
		cacheUntil := time.Now().UTC().Add(duration).Format(http.TimeFormat)
		response.Header.Set("Expires", cacheUntil)
	}

	responseBytes, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
	}

	log.Print("lruttl set")
	c.cache.Add(s, responseBytes)
}

