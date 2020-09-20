package internal

import (
	"bufio"
	"bytes"
	"encoding/gob"
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

type item struct {
	Expiration int64
	Response   []byte
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

	if time.Now().UnixNano() > item.Expiration {
		log.Print(fmt.Sprintf("error:%s", "expired soft ttl"))
		return nil, false
	}

	resp, err :=http.ReadResponse(bufio.NewReader(bytes.NewReader(item.Response)), nil)
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
		return nil, false
	}

	log.Print("lruttl get")
	return resp, true
}

func (c *lruttlCache) Set(s string, response *http.Response, duration time.Duration) {
	responseBytes, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
	}

	var itemBytes bytes.Buffer
	enc := gob.NewEncoder(&itemBytes)
	err = enc.Encode(item{Expiration: time.Now().Add(duration).UnixNano(), Response: responseBytes})
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
	}

	log.Print("lruttl set")
	c.cache.Add(s, item{Expiration: time.Now().Add(duration).UnixNano(), Response: responseBytes})
}
