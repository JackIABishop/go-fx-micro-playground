package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/JackIABishop/go-fx-micro-playground/internal/logging"
)

func validateRates(rates map[string]map[string]float64) error {
	if len(rates) == 0 {
		return errors.New("rates payload is empty")
	}
	for base, targets := range rates {
		if base == "" {
			return errors.New("base currency code cannot be empty")
		}
		if len(targets) == 0 {
			return fmt.Errorf("no target rates provided for %s", base)
		}
		for to, rate := range targets {
			if to == "" {
				return errors.New("target currency code cannot be empty")
			}
			if rate <= 0 {
				return fmt.Errorf("invalid rate for %s->%s: %f", base, to, rate)
			}
		}
	}
	return nil
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Println("💓 /health hit")
	fmt.Fprintln(w, "✅ Rates service is up")
}

func handleRates(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		logging.Logger.Println("📊 GET /rates hit")
		w.Header().Set("Content-Type", "application/json")
		rates := loadRates()
		json.NewEncoder(w).Encode(rates)

	case http.MethodPost:
		logging.Logger.Println("📥 POST /rates hit")
		var newRates map[string]map[string]float64
		if err := json.NewDecoder(r.Body).Decode(&newRates); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		if err := validateRates(newRates); err != nil {
			http.Error(w, fmt.Sprintf("invalid rate data: %v", err), http.StatusBadRequest)
			return
		}
		// Load existing rates from file
		existing := loadRates()
		// Merge new rates into existing
		for base, targets := range newRates {
			if existing[base] == nil {
				existing[base] = targets
			} else {
				for to, rate := range targets {
					existing[base][to] = rate
				}
			}
		}
		// Use merged map for saving
		newRates = existing
		saveRatesToFile(savedRatesFile, newRates)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "rates updated"})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
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
