package internal

import (
	cacheable "github.com/TaylorOno/http-cacheable"
	"net/http"
	"time"
)

type multiStageCache struct {
	l1Cache cacheable.HTTPCacheProvider
	l2Cache cacheable.HTTPCacheProvider
}

func NewMultiStageCache(l1 cacheable.HTTPCacheProvider, l2 cacheable.HTTPCacheProvider) *multiStageCache {
	return &multiStageCache{l1Cache: l1, l2Cache: l2}
}

func (m *multiStageCache) Get(s string) (*http.Response, bool) {
	result, ok := m.l1Cache.Get(s)
	if ok {
		return result, ok
	}

	result, ok = m.l2Cache.Get(s)
	if ok {
		return result, ok
	}

	return nil, false
}

func (m *multiStageCache) Set(s string, response *http.Response, duration time.Duration) {
	m.l1Cache.Set(s,response, duration)
	m.l2Cache.Set(s,response, duration)
}