package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getRates() map[string]map[string]float64 {
	return map[string]map[string]float64{
		"USD": {
			"EUR": 0.92,
			"GBP": 0.78,
			"JPY": 135.33,
		},
		"EUR": {
			"USD": 1.09,
			"GBP": 0.85,
		},
		"GBP": {
			"USD": 1.29,
			"EUR": 1.17,
		},
	}
}

func main() {
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "✅ Rates service is up")
	})

	// Returns a nested JSON structure keyed by base currency
	http.HandleFunc("/rates", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		rates := getRates()
		json.NewEncoder(w).Encode(rates)
	})

	fmt.Println("🚀 Rates service running on :8081")
	// Start HTTP server on port 8081
	http.ListenAndServe(":8081", nil)
}
