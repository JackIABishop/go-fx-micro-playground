package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JackIABishop/go-fx-micro-playground/internal/logging"
)

// TestHandleHealth verifies the /health endpoint responds with 200 OK and correct body.
func TestHandleHealth(t *testing.T) {
	logging.Init()
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

// TestHandleConvert runs multiple scenarios against the /convert endpoint.
func TestHandleConvert(t *testing.T) {
	logging.Init()
	cases := []struct {
		name           string
		query          string
		setupMock      func()
		wantStatus     int
		wantBodyPrefix string
	}{
		{
			name:       "BadAmount",
			query:      "/convert?from=USD&to=EUR&amount=foo",
			setupMock:  func() {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "NoRatesService",
			query: "/convert?from=USD&to=EUR&amount=100",
			setupMock: func() {
				ratesServiceURL = "http://localhost:0/rates"
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:  "ValidConversion",
			query: "/convert?from=USD&to=EUR&amount=100",
			setupMock: func() {
				mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"USD":{"EUR":0.85}}`))
				}))
				ratesServiceURL = mock.URL + "/rates"
				t.Cleanup(func() {
					ratesServiceURL = "http://localhost:8081/rates"
					mock.Close()
				})
			},
			wantStatus:     http.StatusOK,
			wantBodyPrefix: `{"from":"USD","to":"EUR","rate":0.85,"converted":85}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// reset URL
			ratesServiceURL = "http://localhost:8081/rates"
			tc.setupMock()

			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			rr := httptest.NewRecorder()
			handleConvert(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("case %s: expected status %d, got %d", tc.name, tc.wantStatus, rr.Code)
			}
			if tc.name == "ValidConversion" {
				var resp struct {
					From      string  `json:"from"`
					To        string  `json:"to"`
					Rate      float64 `json:"rate"`
					Amount    float64 `json:"amount"`
					Converted float64 `json:"converted"`
				}
				if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
					t.Fatalf("case %s: invalid JSON: %v", tc.name, err)
				}
				if resp.From != "USD" || resp.To != "EUR" || resp.Rate != 0.85 || resp.Amount != 100 || resp.Converted != 85 {
					t.Errorf("case %s: unexpected response: %+v", tc.name, resp)
				}
			} else if tc.wantBodyPrefix != "" {
				if !strings.HasPrefix(rr.Body.String(), tc.wantBodyPrefix) {
					t.Errorf("case %s: unexpected body %q", tc.name, rr.Body.String())
				}
			}
		})
	}
}
