package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JackIABishop/go-fx-micro-playground/internal/logging"
)

func init() {
	logging.Init()
}

// TestHandleRates_Post_Merge verifies POST /rates merges incoming rates with existing ones
func TestHandleRates_Post_Merge(t *testing.T) {
	dir := t.TempDir()
	// Override file paths
	savedRatesFile = filepath.Join(dir, "saved_rates.json")
	newRatesFile = filepath.Join(dir, "new_rates.json")

	// Write initial saved rates: USD->EUR=0.92, GBP->USD=1.27
	initial := `{"USD":{"EUR":0.92},"GBP":{"USD":1.27}}`
	if err := os.WriteFile(savedRatesFile, []byte(initial), 0644); err != nil {
		t.Fatalf("failed to write initial saved rates: %v", err)
	}

	// POST payload updates USD->EUR to 0.95 and adds AUD->USD=0.66
	payload := `{"USD":{"EUR":0.95},"AUD":{"USD":0.66}}`
	req := httptest.NewRequest(http.MethodPost, "/rates", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	logging.Init()
	handleRates(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	// Read back merged rates
	data, err := os.ReadFile(savedRatesFile)
	if err != nil {
		t.Fatalf("failed to read merged rates file: %v", err)
	}
	var merged map[string]map[string]float64
	if err := json.Unmarshal(data, &merged); err != nil {
		t.Fatalf("invalid JSON in merged file: %v", err)
	}
	// Check that USD->EUR was updated
	if merged["USD"]["EUR"] != 0.95 {
		t.Errorf("expected USD->EUR=0.95, got %f", merged["USD"]["EUR"])
	}
	// Check that GBP->USD remains unchanged
	if merged["GBP"]["USD"] != 1.27 {
		t.Errorf("expected GBP->USD=1.27, got %f", merged["GBP"]["USD"])
	}
	// Check that AUD->USD was added
	if merged["AUD"]["USD"] != 0.66 {
		t.Errorf("expected AUD->USD=0.66, got %f", merged["AUD"]["USD"])
	}
}

// TestHandleRates_Post_InvalidJSON verifies POST /rates with malformed JSON returns 400 Bad Request.
func TestHandleRates_Post_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	savedRatesFile = filepath.Join(dir, "saved_rates.json")
	newRatesFile = filepath.Join(dir, "new_rates.json")

	logging.Init()
	req := httptest.NewRequest(http.MethodPost, "/rates", strings.NewReader("{invalid-json}"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handleRates(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for malformed JSON, got %d", rr.Code)
	}
}

// TestHandleRates_Post_MethodNotAllowed verifies unsupported methods on /rates return 405 Method Not Allowed.
func TestHandleRates_Post_MethodNotAllowed(t *testing.T) {
	logging.Init()
	req := httptest.NewRequest(http.MethodDelete, "/rates", nil)
	rr := httptest.NewRecorder()

	handleRates(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 for method not allowed, got %d", rr.Code)
	}
}

// TestHandleHealth verifies that the /health endpoint responds with 200 OK and the expected body
func TestHandleHealth(t *testing.T) {
	// Create a GET request for the /health endpoint
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	handleHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if body := rr.Body.String(); body != "✅ Rates service is up\n" {
		t.Errorf("unexpected body: %q", body)
	}
}

// TestHandleRates verifies that the /rates endpoint returns valid JSON containing currency rate data
func TestHandleRates(t *testing.T) {
	dir := t.TempDir()
	savedRatesFile = filepath.Join(dir, "saved_rates.json")
	newRatesFile = filepath.Join(dir, "new_rates.json")

	sample := `{"USD":{"EUR":0.92}}`
	if err := os.WriteFile(savedRatesFile, []byte(sample), 0644); err != nil {
		t.Fatalf("failed to write saved rates file: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/rates", nil)
	rr := httptest.NewRecorder()
	handleRates(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected JSON content-type, got %q", ct)
	}
	var resp map[string]map[string]float64
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if _, ok := resp["USD"]["EUR"]; !ok {
		t.Errorf("expected USD to EUR rate in response, got %v", resp)
	}
}

// TestHandleRates_Post verifies POST /rates updates the saved file.
func TestHandleRates_Post(t *testing.T) {
	dir := t.TempDir()
	savedRatesFile = filepath.Join(dir, "saved_rates.json")
	newRatesFile = filepath.Join(dir, "new_rates.json")

	payload := `{"USD":{"GBP":0.78}}`
	req := httptest.NewRequest(http.MethodPost, "/rates", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handleRates(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	data, err := os.ReadFile(savedRatesFile)
	if err != nil {
		t.Fatalf("failed to read saved rates file: %v", err)
	}
	var resp map[string]map[string]float64
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("invalid JSON in saved file: %v", err)
	}
	if resp["USD"]["GBP"] != 0.78 {
		t.Errorf("expected persisted rate 0.78, got %f", resp["USD"]["GBP"])
	}
}
