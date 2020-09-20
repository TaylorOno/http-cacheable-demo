package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
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

	item, ok := value.([]byte)
	if !ok {
		return nil, false
	}

	resp, err :=http.ReadResponse(bufio.NewReader(bytes.NewReader(item)), nil)
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
		return nil, false
	}

	log.Print("gocahe get")
	return resp, true
}

func (c *ttlCache) Set(s string, response *http.Response, duration time.Duration) {
	responseBytes, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
	}

	log.Print("gocahe set")
	c.cache.Set(s, responseBytes, duration)
}
