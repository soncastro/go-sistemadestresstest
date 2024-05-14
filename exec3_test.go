package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestTotalRequests(t *testing.T) {
	var totalRequests int
	var mtx sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		mtx.Lock()
		totalRequests++
		mtx.Unlock()
	}))
	defer server.Close()

	expectedRequests := 10
	sem := make(chan bool, 2)
	var wg sync.WaitGroup

	for i := 0; i < expectedRequests; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			sem <- true
			_, err := http.Get(server.URL)

			if err != nil {
				t.Errorf("Request failed: %v", err)
			}

			<-sem
		}()
	}

	wg.Wait()

	if totalRequests != expectedRequests {
		t.Errorf("Expected %v requests but got %v", expectedRequests, totalRequests)
	}
}
