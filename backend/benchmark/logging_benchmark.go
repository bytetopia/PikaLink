package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"pikalink-backend/logging"
	"time"
)

func benchmarkLogging() {
	// Create a test HTTP request
	req := httptest.NewRequest("GET", "/test-url", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	req.Header.Set("Referer", "https://example.com")

	const numRequests = 1000

	fmt.Printf("Benchmarking logging performance with %d requests...\n\n", numRequests)

	// Benchmark synchronous logging
	fmt.Println("Testing SYNCHRONOUS logging (current implementation):")
	start := time.Now()
	for i := 0; i < numRequests; i++ {
		err := logging.LogLinkAccess(req, "test123", "https://example.com", http.StatusFound)
		if err != nil {
			fmt.Printf("Error in sync logging: %v\n", err)
		}
	}
	syncDuration := time.Since(start)
	fmt.Printf("Synchronous logging took: %v\n", syncDuration)
	fmt.Printf("Average per request: %v\n", syncDuration/numRequests)
	fmt.Printf("Requests per second: %.2f\n\n", float64(numRequests)/syncDuration.Seconds())

	// Benchmark asynchronous logging
	fmt.Println("Testing ASYNCHRONOUS logging (new implementation):")
	start = time.Now()
	for i := 0; i < numRequests; i++ {
		logging.LogLinkAccessAsync(req, "test123", "https://example.com", http.StatusFound)
	}
	asyncDuration := time.Since(start)
	fmt.Printf("Asynchronous logging took: %v\n", asyncDuration)
	fmt.Printf("Average per request: %v\n", asyncDuration/numRequests)
	fmt.Printf("Requests per second: %.2f\n\n", float64(numRequests)/asyncDuration.Seconds())

	// Performance improvement
	improvement := float64(syncDuration) / float64(asyncDuration)
	fmt.Printf("Performance improvement: %.2fx faster\n", improvement)
	fmt.Printf("Latency reduction: %v per request\n", (syncDuration-asyncDuration)/numRequests)

	// Wait a bit for async processing to complete
	fmt.Println("\nWaiting 10 seconds for async processing to complete...")
	time.Sleep(10 * time.Second)

	// Stop the async logger to flush remaining entries
	logger := logging.GetAsyncLogger()
	logger.Stop()
}

func main() {
	benchmarkLogging()
}