package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "✅ Rates service is up")
	})

	// Returns a hardcoded exchange rate (USD to EUR)
	http.HandleFunc("/rates", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Simulate a basic exchange rate as JSON
		w.Write([]byte(`{"USD_EUR": 0.92}`))
	})

	fmt.Println("🚀 Rates service running on :8081")
	// Start HTTP server on port 8081
	http.ListenAndServe(":8081", nil)
}
