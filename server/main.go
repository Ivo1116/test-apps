package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Create a struct to simulate the large data stream
type LargeDataStream struct {
	dataSize  int64
	chunkSize int64
	bytesSent int64
}

// Constructor to initialize LargeDataStream
func NewLargeDataStream() *LargeDataStream {
	return &LargeDataStream{
		dataSize:  20 * 1024 * 1024 * 1024, // 20GB
		chunkSize: 1024 * 1024 * 10,        // 10MB per chunk
		bytesSent: 0,
	}
}

// Read method to simulate sending 10MB chunks at a time
func (lds *LargeDataStream) Read(p []byte) (n int, err error) {
	if lds.bytesSent < lds.dataSize {
		// Fill the buffer with dummy data (e.g., 'A')
		for i := range p {
			p[i] = 'A'
		}
		lds.bytesSent += lds.chunkSize
		if lds.bytesSent > lds.dataSize {
			n = int(lds.dataSize - (lds.bytesSent - lds.chunkSize))
			err = io.EOF
		} else {
			n = len(p)
		}
	} else {
		err = io.EOF
	}
	return
}

// Ping handler that returns a timestamp
func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	timestamp := map[string]string{"timestamp": time.Now().Format(time.RFC3339)}
	// Respond with the timestamp in JSON format
	if err := json.NewEncoder(w).Encode(timestamp); err != nil {
		http.Error(w, "Failed to encode timestamp", http.StatusInternalServerError)
	}
}

// Download handler that streams the large data
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for the large file download
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", 20*1024*1024*1024)) // 20GB
	w.Header().Set("Connection", "keep-alive")                             // Keep the connection alive

	// Create a new LargeDataStream instance
	largeDataStream := NewLargeDataStream()

	// Create a buffer large enough to hold a 10MB chunk
	buffer := make([]byte, largeDataStream.chunkSize)

	// Copy data from the LargeDataStream to the response
	for {
		// Read the next chunk of data
		n, err := largeDataStream.Read(buffer)
		if err != nil && err != io.EOF {
			http.Error(w, "Failed to read data", http.StatusInternalServerError)
			return
		}
		if n > 0 {
			// Write the chunk to the response
			if _, err := w.Write(buffer[:n]); err != nil {
				http.Error(w, "Failed to send data", http.StatusInternalServerError)
				return
			}
		}
		if err == io.EOF {
			break
		}
	}
}

func main() {
	// Define the server and routes
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/download", downloadHandler)

	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)

	// Start the server
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
