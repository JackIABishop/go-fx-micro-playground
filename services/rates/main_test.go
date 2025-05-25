package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JackIABishop/go-fx-micro-playground/internal/logging"
)

func init() {
	logging.Init()
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
	// Create a GET request for the /rates endpoint
	req := httptest.NewRequest(http.MethodGet, "/rates", nil)
	rr := httptest.NewRecorder()
	handleRates(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected JSON content-type, got %q", ct)
	}
	var rates map[string]map[string]float64
	// Parse JSON response into a nested map[string]map[string]float64
	if err := json.Unmarshal(rr.Body.Bytes(), &rates); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	// Assert that the USD to EUR conversion rate is present
	if _, ok := rates["USD"]["EUR"]; !ok {
		t.Errorf("expected USD to EUR rate in response, got %v", rates)
	}
}
