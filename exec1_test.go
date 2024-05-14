package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(rw, "Dummy server response")
	}))
	defer server.Close()

	url := server.URL // O servidor fornece a URL para o teste
	_, err := http.Get(url)
	if err != nil {
		t.Errorf("TestHttpRequest failed: %v", err)
	}
}
