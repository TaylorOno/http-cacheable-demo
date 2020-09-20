package internal

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
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

	var item item
	dec := gob.NewDecoder(bytes.NewReader(memCacheItem.Value))
	err = dec.Decode(&item)
	if err  != nil  {
		log.Print(fmt.Sprintf("error:%s", err))
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

	log.Print("memcache get")
	return resp, true
}

func (m *memCacheClient) Set(s string, response *http.Response, duration time.Duration) {
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

	err = m.client.Set(&memcache.Item{Key: s, Value: itemBytes.Bytes()})
	if err != nil {
		log.Print(fmt.Sprintf("error:%s", err))
	}
	log.Print("memcache set")
}