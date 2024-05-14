package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"time"
)

type TestReport struct {
	TotalRequests      int
	SuccessfulRequests int
	FailedRequests     int
	ResponseCodes      map[int]int
	TotalTime          time.Duration
}

func main() {
	url := flag.String("url", "http://localhost:8080", "URL to test.")
	requests := flag.Int("requests", 50, "Number of requests to make.")
	concurrency := flag.Int("concurrency", 10, "Number of concurrent requests.")
	flag.Parse()

	report := TestReport{
		ResponseCodes: make(map[int]int),
	}

	sem := make(chan bool, *concurrency)
	var wg sync.WaitGroup
	var m sync.Mutex

	startTime := time.Now()

	for i := 0; i < *requests; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			sem <- true
			response, err := http.Get(*url)

			if err != nil {
				log.Printf("Request failed: %v\n", err)
				m.Lock()
				report.FailedRequests++
				report.ResponseCodes[0]++ // Use 0 for failed requests
				m.Unlock()
			} else {
				m.Lock()
				report.ResponseCodes[response.StatusCode]++
				if response.StatusCode == http.StatusOK {
					report.SuccessfulRequests++
				}
				m.Unlock()
			}

			<-sem
		}()
	}

	wg.Wait()

	totalTime := time.Since(startTime)
	report.TotalTime = totalTime
	report.TotalRequests = *requests

	printReport(report)
}

func printReport(report TestReport) {
	log.Printf("Total time spent: %v", report.TotalTime)
	log.Printf("Total requests: %v", report.TotalRequests)
	log.Printf("Successful requests: %v", report.SuccessfulRequests)
	//log.Printf("Failed requests: %v", report.FailedRequests)

	for statusCode, count := range report.ResponseCodes {
		if statusCode == 0 {
			log.Printf("Failed requests: %v", count)
		} else {
			log.Printf("HTTP %v: %v", statusCode, count)
		}
	}
}
