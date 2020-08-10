package routes

import (
	"compress/gzip"
	"io/ioutil"
	"log"
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
	if resp.Uncompressed {
		defer resp.Body.Close()
		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if !resp.Uncompressed {
		zipReader, err := gzip.NewReader(resp.Body)
		defer zipReader.Close()
		if err != nil {
			log.Fatal(err)
		}
		result, err = ioutil.ReadAll(zipReader)
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
	if resp.Uncompressed {
		defer resp.Body.Close()
		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if !resp.Uncompressed {
		zipReader, err := gzip.NewReader(resp.Body)
		defer zipReader.Close()
		if err != nil {
			log.Fatal(err)
		}
		result, err = ioutil.ReadAll(zipReader)
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
	if resp.Uncompressed {
		defer resp.Body.Close()
		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if !resp.Uncompressed {
		zipReader, err := gzip.NewReader(resp.Body)
		defer zipReader.Close()
		if err != nil {
			log.Fatal(err)
		}
		result, err = ioutil.ReadAll(zipReader)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(result)
	return
}
