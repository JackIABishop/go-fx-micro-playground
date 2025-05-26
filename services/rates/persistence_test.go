package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Test readRatesFromFile with valid data
func TestReadRatesFromFile_Valid(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "valid_rates.json")
	data := `{
		"USD": {"EUR": 0.92, "GBP": 0.78},
		"EUR": {"USD": 1.09}
	}`
	if err := os.WriteFile(filePath, []byte(data), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	rates, err := readRatesFromFile(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rates["USD"]["EUR"] != 0.92 {
		t.Errorf("expected 0.92, got %f", rates["USD"]["EUR"])
	}
}

// Test readRatesFromFile with malformed JSON
func TestReadRatesFromFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "invalid_rates.json")
	data := `{"USD": {"EUR": "oops"}}`
	if err := os.WriteFile(filePath, []byte(data), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := readRatesFromFile(filePath)
	if err == nil {
		t.Error("expected error for malformed JSON, got nil")
	}
}

// Test saveRatesToFile and read back to verify
func TestSaveAndReadRatesToFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "saved_rates.json")
	original := map[string]map[string]float64{
		"USD": {"EUR": 0.92},
		"EUR": {"USD": 1.09},
	}

	saveRatesToFile(filePath, original)

	readBack, err := readRatesFromFile(filePath)
	if err != nil {
		t.Fatalf("unexpected error reading back saved rates: %v", err)
	}
	if readBack["USD"]["EUR"] != 0.92 {
		t.Errorf("expected 0.92, got %f", readBack["USD"]["EUR"])
	}
}

func TestLoadRates_FallbackToSaved(t *testing.T) {
	dir := t.TempDir()

	// Override file paths
	newRatesFile = filepath.Join(dir, "new_rates.json")
	savedRatesFile = filepath.Join(dir, "saved_rates.json")

	// Write only saved file with valid data
	savedData := `{"USD": {"EUR": 0.92}}`
	if err := os.WriteFile(savedRatesFile, []byte(savedData), 0644); err != nil {
		t.Fatalf("failed to write saved file: %v", err)
	}

	rates := loadRates()
	if rates["USD"]["EUR"] != 0.92 {
		t.Errorf("expected saved rate 0.92, got %f", rates["USD"]["EUR"])
	}
}

func TestLoadRates_FallbackToEmpty(t *testing.T) {
	dir := t.TempDir()

	// Override file paths to missing files
	newRatesFile = filepath.Join(dir, "new_rates.json")
	savedRatesFile = filepath.Join(dir, "saved_rates.json")

	rates := loadRates()
	if len(rates) != 0 {
		t.Errorf("expected empty rates, got %+v", rates)
	}
}
