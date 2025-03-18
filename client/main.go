package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Measure the ping to the server
func pingServer(url string) (time.Duration, error) {
	start := time.Now()
	_, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	return time.Since(start), nil
}

// Continuously ping the server and measure ping time with a specified interval
func pingWithInterval(url string, pingInterval time.Duration) {
	for {
		pingDuration, err := pingServer(url + "/ping")
		if err != nil {
			log.Printf("Error pinging server: %v", err)
			continue
		}
		fmt.Printf("Ping time: %v\n", pingDuration)
		time.Sleep(pingInterval)
	}
}

// Main function
func main() {
	// Command-line flags for URL and ping interval
	serverURL := flag.String("url", "http://localhost:8080", "Server URL")
	pingInterval := flag.Int("interval", 5, "Ping interval in seconds")
	flag.Parse()

	// Convert ping interval from seconds to time.Duration
	interval := time.Duration(*pingInterval) * time.Second

	// Start pinging the server with the specified parameters
	pingWithInterval(*serverURL, interval)
}
