package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	// Health check for the gateway service
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "✅ Gateway is up")
	})

	// Endpoint to perform a currency conversion
	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		// Call the /rates endpoint from the rates service
		resp, err := http.Get("http://localhost:8081/rates")
		if err != nil {
			http.Error(w, "Failed to contact rates service", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Decode the JSON response into a map
		var rates map[string]float64
		if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
			http.Error(w, "Bad response from rates service", http.StatusInternalServerError)
			return
		}

		// Use the USD_EUR rate to simulate converting 100 USD to EUR
		usdToEur := rates["USD_EUR"]
		fmt.Fprintf(w, "💱 100 USD is %.2f EUR\n", 100*usdToEur)
	})

	fmt.Println("🚀 Gateway running on :8080")
	// Start HTTP server on port 8080
	http.ListenAndServe(":8080", nil)
}
