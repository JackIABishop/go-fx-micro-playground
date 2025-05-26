package main

import (
	"encoding/json"
	"os"

	"github.com/JackIABishop/go-fx-micro-playground/internal/logging"
)

var savedRatesFile = "/app/services/rates/saved_rates.json"
var newRatesFile = "/app/services/rates/new_rates.json"

func readRatesFromFile(path string) (map[string]map[string]float64, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		logging.Logger.Printf("⚠️ Failed to load rates from %s: %v", path, err)
		return nil, err
	}
	var rates map[string]map[string]float64
	if err := json.Unmarshal(file, &rates); err != nil {
		logging.Logger.Printf("⚠️ Failed to load rates from %s: %v", path, err)
		return nil, err
	}
	logging.Logger.Printf("✅ Loaded rates from %s", path)
	return rates, nil
}

func loadRates() map[string]map[string]float64 {
	cwd, _ := os.Getwd()
	logging.Logger.Printf("🐛 Current working directory: %s", cwd)
	rates, err := readRatesFromFile(newRatesFile)
	if err == nil {
		saveRatesToFile(savedRatesFile, rates)
		return rates
	}

	rates, err = readRatesFromFile(savedRatesFile)
	if err == nil {
		return rates
	}

	// Fallback to empty set
	logging.Logger.Printf("❌ No rates available from any source, returning empty set")
	return map[string]map[string]float64{}
}

// TODO: Use `saved_rates.json` as persistent cache for fetched rates.
// Fallback to default rates if not available. Will also support syncing with an external API.

func saveRatesToFile(path string, rates map[string]map[string]float64) {
	data, err := json.MarshalIndent(rates, "", "  ")
	if err != nil {
		logging.Logger.Printf("❌ Failed to marshal rates: %v", err)
		return
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		logging.Logger.Printf("❌ Failed to write rates to file %s: %v", path, err)
	} else {
		logging.Logger.Printf("💾 Rates successfully saved to %s", path)
	}
}
