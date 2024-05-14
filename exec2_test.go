package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

var maxGoroutines int
var currentGoroutines int

func TestConcurrency(t *testing.T) {
	maxGoroutines = 0
	currentGoroutines = 0
	mtx := &sync.Mutex{}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Atualiza o contador de goroutines ativas
		mtx.Lock()
		currentGoroutines++
		if currentGoroutines > maxGoroutines {
			maxGoroutines = currentGoroutines
		}
		mtx.Unlock()

		// Simula um tempo de processamento para permitir que outras goroutines sejam executadas
		time.Sleep(100 * time.Millisecond)

		mtx.Lock()
		currentGoroutines--
		mtx.Unlock()
	}))

	defer server.Close()

	concurrency := 5
	sem := make(chan bool, concurrency)
	var wg sync.WaitGroup

	for i := 0; i < 15; i++ {
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

	if maxGoroutines != concurrency {
		t.Errorf("Expected %v maximum goroutines but got %v", concurrency, maxGoroutines)
	}
}
