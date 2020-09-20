package routes

import (
	"github.com/TaylorOno/http-cacheable-demo/internal"
	"net/http"
	"strings"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Server struct {
	GoCacheClient     Client
	LRUCacheClient    Client
	LRUTTLCacheClient Client
	MemCacheClient    Client
	MultiStageClient  Client
}

func (s *Server) GoCache(w http.ResponseWriter, req *http.Request) {
	proxyRequest, _ := http.NewRequest(req.Method, "https://reddit.com", req.Body)
	proxyRequest.Header = req.Header
	proxyRequest.URL.Path = strings.Replace(req.URL.Path, "/goCache", "", -1)
	proxyRequest.URL.RawQuery = req.URL.RawQuery
	resp, err := s.GoCacheClient.Do(proxyRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var result []byte
	if resp.Body != nil {
		result, err = internal.ReadBody(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(result)
	return
}

func (s *Server) LRUCache(w http.ResponseWriter, req *http.Request) {
	proxyRequest, _ := http.NewRequest(req.Method, "https://reddit.com", req.Body)
	proxyRequest.Header = req.Header
	proxyRequest.URL.Path = strings.Replace(req.URL.Path, "/lruCache", "", -1)
	proxyRequest.URL.RawQuery = req.URL.RawQuery
	resp, err := s.LRUCacheClient.Do(proxyRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var result []byte
	if resp.Body != nil {
		result, err = internal.ReadBody(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(result)
	return
}

func (s *Server) LRUTTLCache(w http.ResponseWriter, req *http.Request) {
	proxyRequest, _ := http.NewRequest(req.Method, "https://reddit.com", req.Body)
	proxyRequest.Header = req.Header
	proxyRequest.URL.Path = strings.Replace(req.URL.Path, "/lruttlCache", "", -1)
	proxyRequest.URL.RawQuery = req.URL.RawQuery
	resp, err := s.LRUTTLCacheClient.Do(proxyRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var result []byte
	if resp.Body != nil {
		result, err = internal.ReadBody(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(result)
	return
}

func (s *Server) MemCache(w http.ResponseWriter, req *http.Request) {
	proxyRequest, _ := http.NewRequest(req.Method, "https://reddit.com", req.Body)
	proxyRequest.Header = req.Header
	proxyRequest.URL.Path = strings.Replace(req.URL.Path, "/memCached", "", -1)
	proxyRequest.URL.RawQuery = req.URL.RawQuery
	resp, err := s.MemCacheClient.Do(proxyRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var result []byte
	if resp.Body != nil {
		result, err = internal.ReadBody(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(result)
	return
}

func (s *Server) MultiStageCache(w http.ResponseWriter, req *http.Request) {
	proxyRequest, _ := http.NewRequest(req.Method, "https://reddit.com", req.Body)
	proxyRequest.Header = req.Header
	proxyRequest.URL.Path = strings.Replace(req.URL.Path, "/multiStage", "", -1)
	proxyRequest.URL.RawQuery = req.URL.RawQuery
	resp, err := s.MultiStageClient.Do(proxyRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var result []byte
	if resp.Body != nil {
		result, err = internal.ReadBody(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(result)
	return
}
