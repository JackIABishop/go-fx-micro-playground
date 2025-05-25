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

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "✅ Rates service is up")
}

func handleRates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rates := getRates()
	json.NewEncoder(w).Encode(rates)
}

// ⚠️ NOTE:
// This is a nice clean route setup for learning purposes — easy to read and reason about.
// But in production, you'd use a router like chi or echo to handle things like:
// - Middleware (auth, logging, CORS, etc)
// - Cleaner route grouping and param handling
// - Better testability and DI-friendly structure
// I'm keeping it simple for now while I am focused on the basics.
func setupRoutes() {
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/rates", handleRates)
}

func main() {
	setupRoutes()
	fmt.Println("🚀 Rates service running on :8081")
	// Start HTTP server on port 8081
	http.ListenAndServe(":8081", nil)
}
