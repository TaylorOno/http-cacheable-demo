package ttl

import (
	"bytes"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"net/http"
	"time"
)

type ttlCache struct {
	cache *cache.Cache
}

type item struct {
	uncompressed bool
	statusCode   int
	body         []byte
}

func NewTTLCache() *ttlCache {
	return &ttlCache{cache: cache.New(5*time.Minute, 10*time.Minute)}
}

func (c *ttlCache) Get(s string) (*http.Response, bool) {
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

func (c *ttlCache) Set(s string, response *http.Response, duration time.Duration) {
	var bodyBytes []byte
	if response.Body != nil {
		defer response.Body.Close()
		bodyBytes, _ = ioutil.ReadAll(response.Body)
	}
	c.cache.Set(s, item{uncompressed: response.Uncompressed, statusCode: response.StatusCode, body: bodyBytes}, duration)
	response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
}
