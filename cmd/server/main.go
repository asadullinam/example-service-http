package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type healthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(responseWriter http.ResponseWriter, _ *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(responseWriter).Encode(map[string]string{
			"message": "example-service is running",
		})
	})
	mux.HandleFunc("/health", func(responseWriter http.ResponseWriter, _ *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(responseWriter).Encode(healthResponse{
			Status:    "ok",
			Service:   "example-service",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	})

	address := ":" + port
	log.Printf("example-service started on %s", address)
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
