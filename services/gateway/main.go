package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "✅ Gateway is up")
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	// Call the /rates endpoint from the rates service
	resp, err := http.Get("http://localhost:8081/rates")
	if err != nil {
		http.Error(w, "❌ Failed to contact rates service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode the JSON response into a nested map
	var rates map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		http.Error(w, "❌ Bad response from rates service", http.StatusInternalServerError)
		return
	}

	// Simulate converting 100 units from various currencies
	baseAmount := 100.0
	output := "💸 Currency Conversions for 100 units:\n\n"
	for base, conversions := range rates {
		for target, rate := range conversions {
			converted := baseAmount * rate
			output += fmt.Sprintf("🔁 %3.0f %s ➡️ %.2f %s\n", baseAmount, base, converted, target)
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(output))
}

// setupRoutes registers HTTP handlers for the gateway service.
// This is a simplified setup for learning purposes.
// In production systems, a router would typically be used for more flexibility and features.
func setupRoutes() {
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/convert", handleConvert)
}

func main() {
	setupRoutes()

	fmt.Println("🚀 Gateway running on :8080")
	// Start HTTP server on port 8080
	http.ListenAndServe(":8080", nil)
}
