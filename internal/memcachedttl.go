package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type memCacheClient struct {
	client *memcache.Client
}

func NewMemCacheClient() *memCacheClient {
	return &memCacheClient{
		client: memcache.New("localhost:11211"),
	}
}

func (m *memCacheClient) Get(s string) (*http.Response, bool) {
	memCacheItem, err := m.client.Get(s)
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
		return nil, false
	}

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(memCacheItem.Value)), nil)
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
		return nil, false
	}

	expires, err := http.ParseTime(resp.Header.Get("Expires"))
	if err != nil || time.Now().UTC().After(expires) {
		log.Print(fmt.Sprintf("error:%s", "expired soft ttl"))
		return resp, false
	}

	log.Print("memcached get")
	return resp, true
}

func (m *memCacheClient) Set(s string, response *http.Response, duration time.Duration) {
	if duration > 0 {
		cacheUntil := time.Now().UTC().Add(duration).Format(http.TimeFormat)
		response.Header.Set("Expires", cacheUntil)
	}

	responseBytes, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
	}

	err = m.client.Set(&memcache.Item{Key: s, Value: responseBytes})
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
	}
	log.Print("memcached set")
}
