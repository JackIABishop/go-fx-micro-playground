package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JackIABishop/go-fx-micro-playground/internal/logging"
)

func handleHealth(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Println("💓 /health hit")
	fmt.Fprintln(w, "✅ Rates service is up")
}

func handleRates(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Println("📊 /rates hit")
	w.Header().Set("Content-Type", "application/json")
	rates := loadRates()
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
	logging.Init()
	setupRoutes()
	logging.Logger.Println("🚀 Rates service running on :8081")
	// Start HTTP server on port 8081
	http.ListenAndServe(":8081", nil)
}
