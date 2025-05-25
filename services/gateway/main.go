package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/JackIABishop/go-fx-micro-playground/internal/logging"
)

func handleHealth(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Println("💓 /health hit")
	fmt.Fprintln(w, "✅ Gateway is up")
}

func handleConvert(w http.ResponseWriter, r *http.Request) {

	// Parse query parameters from the URL: from, to, and amount
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amountStr := r.URL.Query().Get("amount")

	// Convert the amount parameter from string to float64
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "❌ Invalid amount parameter", http.StatusBadRequest)
		return
	}

	logging.Logger.Printf("💬 Received conversion request: from=%s to=%s amount=%f", from, to, amount)

	// Call the /rates endpoint from the rates service to get current currency rates
	resp, err := http.Get("http://localhost:8081/rates")
	if err != nil {
		logging.Logger.Printf("❌ Error contacting rates service: %v", err)
		http.Error(w, "❌ Failed to contact rates service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode the JSON response from the rates service into a nested map
	var rates map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		logging.Logger.Printf("❌ Error decoding rates response: %v", err)
		http.Error(w, "❌ Bad response from rates service", http.StatusInternalServerError)
		return
	}

	// Look up the map of target currency rates for the 'from' currency
	rateMap, ok := rates[from]
	if !ok {
		logging.Logger.Printf("❌ Unsupported base currency: %s", from)
		http.Error(w, "❌ Unsupported currency: "+from, http.StatusBadRequest)
		return
	}

	// Look up the exchange rate from 'from' to 'to' currency
	targetRate, ok := rateMap[to]
	if !ok {
		logging.Logger.Printf("❌ Unsupported target currency: from=%s to=%s", from, to)
		http.Error(w, "❌ Unsupported target currency: "+to, http.StatusBadRequest)
		return
	}

	// Calculate the converted amount using the exchange rate
	converted := amount * targetRate

	logging.Logger.Printf("✅ Converted %.2f %s to %.2f %s using rate %.4f", amount, from, converted, to, targetRate)

	// Prepare the JSON response with conversion details
	result := map[string]interface{}{
		"from":      from,
		"to":        to,
		"amount":    amount,
		"rate":      targetRate,
		"converted": converted,
	}

	// Set response headers and encode the result as JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(result)
}

// setupRoutes registers HTTP handlers for the gateway service.
// This is a simplified setup for learning purposes.
// In production systems, a router would typically be used for more flexibility and features.
func setupRoutes() {
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/convert", handleConvert)
}

func main() {
	logging.Init()
	logging.Logger.Println("🚀 Gateway running on :8080")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
