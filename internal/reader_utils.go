package internal

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func ReadBody(resp *http.Response) ([]byte, error){
	if !resp.Uncompressed {
		return uncompressAndRead(resp.Body)
	}
	return read(resp.Body)
}

func read(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return ioutil.ReadAll(body)
}

func uncompressAndRead(body io.ReadCloser) ([]byte, error) {
	zipReader, err := gzip.NewReader(body)
	defer zipReader.Close()
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.ReadAll(zipReader)
}
