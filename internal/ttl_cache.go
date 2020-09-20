package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/patrickmn/go-cache"
)

type ttlCache struct {
	cache *cache.Cache
}

func NewTTLCache(duration time.Duration) *ttlCache {
	return &ttlCache{cache: cache.New(duration, 5*time.Minute)}
}

func (c *ttlCache) Get(s string) (*http.Response, bool) {
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

	log.Print("gocache get")
	return resp, true
}

func (c *ttlCache) Set(s string, response *http.Response, duration time.Duration) {
	responseBytes, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
	}

	log.Print("gocache set")
	c.cache.SetDefault(s, responseBytes)
}
