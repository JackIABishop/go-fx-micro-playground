package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JackIABishop/go-fx-micro-playground/internal/logging"
)

// ratesServiceURL is the endpoint used to fetch FX rates; can be overridden in tests.
var ratesServiceURL = "http://localhost:8081/rates"

func init() {
	logging.Init()
}

// TestHandleHealth verifies the /health endpoint responds with 200 OK and correct body.
func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	handleHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if body := rr.Body.String(); body != "✅ Gateway is up\n" {
		t.Errorf("unexpected body: %q", body)
	}
}

// TestHandleConvert_BadAmount verifies invalid amount yields 400 Bad Request.
func TestHandleConvert_BadAmount(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/convert?from=USD&to=EUR&amount=not-a-number", nil)
	rr := httptest.NewRecorder()
	handleConvert(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid amount, got %d", rr.Code)
	}
}

// TestHandleConvert_NoRatesService simulates the rates service being down and expects 500.
func TestHandleConvert_NoRatesService(t *testing.T) {
	originalURL := ratesServiceURL
	ratesServiceURL = "http://localhost:0/rates"
	defer func() { ratesServiceURL = originalURL }()

	req := httptest.NewRequest(http.MethodGet, "/convert?from=USD&to=EUR&amount=100", nil)
	rr := httptest.NewRecorder()
	handleConvert(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 when rates service unavailable, got %d", rr.Code)
	}
}
