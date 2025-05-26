package main

import (
	"encoding/json"
	"os"

	"github.com/JackIABishop/go-fx-micro-playground/internal/logging"
)

var ratesFile = "rates.json"

func loadRates() map[string]map[string]float64 {
	file, err := os.ReadFile(ratesFile)
	if err != nil {
		logging.Logger.Printf("⚠️ Could not read %s: %v — using default rates", ratesFile, err)
		return getDefaultRates()
	}

	var rates map[string]map[string]float64
	if err := json.Unmarshal(file, &rates); err != nil {
		logging.Logger.Printf("⚠️ Invalid JSON in %s: %v — using default rates", ratesFile, err)
		return getDefaultRates()
	}

	logging.Logger.Printf("✅ Loaded rates from %s", ratesFile)
	return rates
}

// TODO: I'd like to save this data in a saved_rates.json file rather than hardcoding it.
// TODO: THEN I can have a function which tries to retrieve the latest rates from the 'api_file' and save it into saved_rates.json
func getDefaultRates() map[string]map[string]float64 {
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
