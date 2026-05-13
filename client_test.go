package pretium

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// rewriteTransport redirects every outbound request to a test server while
// preserving the path, query and body produced by the SDK. This lets us assert
// against the path the SDK builds without exposing the base URL as a var.
type rewriteTransport struct {
	target *url.URL
}

func (rt *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = rt.target.Scheme
	req.URL.Host = rt.target.Host
	return http.DefaultTransport.RoundTrip(req)
}

// newTestClient spins up an httptest server and returns a PretiumClient wired
// to send requests there, plus the server so the caller can close it.
func newTestClient(t *testing.T, handler http.HandlerFunc) (*PretiumClient, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	u, err := url.Parse(ts.URL)
	if err != nil {
		ts.Close()
		t.Fatalf("parse test server URL: %v", err)
	}

	client := New("settlement-addr", "test-api-key", "https://example.com/cb")
	client.SetHTTPClient(&http.Client{Transport: &rewriteTransport{target: u}})
	return client, ts
}

func TestClient_Onramp_RequestShape(t *testing.T) {
	var capturedPath, capturedMethod, capturedAPIKey, capturedContentType string
	var capturedBody map[string]any

	client, ts := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedMethod = r.Method
		capturedAPIKey = r.Header.Get("x-api-key")
		capturedContentType = r.Header.Get("Content-Type")

		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &capturedBody)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"code": 200,
			"message": "ok",
			"data": {"transaction_code":"tx-1","status":"PENDING","message":"queued"}
		}`))
	})
	defer ts.Close()

	resp, err := client.Onramp(context.Background(), KES, OnrampBody{
		Shortcode:     "0712345678",
		Amount:        100,
		MobileNetwork: SAFARICOM,
		Chain:         CELO,
		Asset:         CUSD,
		Address:       "0xabc",
	})
	if err != nil {
		t.Fatalf("Onramp returned error: %v", err)
	}

	if capturedMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", capturedMethod)
	}
	if capturedPath != "/v1/onramp/KES" {
		t.Errorf("path = %q, want /v1/onramp/KES", capturedPath)
	}
	if capturedAPIKey != "test-api-key" {
		t.Errorf("x-api-key = %q, want test-api-key", capturedAPIKey)
	}
	if capturedContentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", capturedContentType)
	}
	if capturedBody["callback_url"] != "https://example.com/cb" {
		t.Errorf("callback_url = %v, want https://example.com/cb", capturedBody["callback_url"])
	}
	if capturedBody["shortcode"] != "0712345678" {
		t.Errorf("shortcode = %v, want 0712345678", capturedBody["shortcode"])
	}
	if capturedBody["asset"] != CUSD {
		t.Errorf("asset = %v, want %s", capturedBody["asset"], CUSD)
	}

	if resp.Data.TransactionCode != "tx-1" {
		t.Errorf("TransactionCode = %q, want tx-1", resp.Data.TransactionCode)
	}
	if resp.Data.Status != StatusPending {
		t.Errorf("Status = %q, want %q", resp.Data.Status, StatusPending)
	}
}

func TestClient_ExchangeRate_RequestAndResponse(t *testing.T) {
	var capturedPath string
	var capturedBody map[string]any

	client, ts := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &capturedBody)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"code": 200,
			"message": "Exchange rates",
			"data": {"buying_rate": 128.15, "selling_rate": 130.75}
		}`))
	})
	defer ts.Close()

	resp, err := client.ExchangeRate(context.Background(), ExchangeRateBody{CurrencyCode: KES})
	if err != nil {
		t.Fatalf("ExchangeRate returned error: %v", err)
	}

	if capturedPath != "/v1/exchange-rate" {
		t.Errorf("path = %q, want /v1/exchange-rate", capturedPath)
	}
	if capturedBody["currency_code"] != KES {
		t.Errorf("currency_code = %v, want %s", capturedBody["currency_code"], KES)
	}
	if resp.Data.BuyingRate != 128.15 {
		t.Errorf("BuyingRate = %v, want 128.15", resp.Data.BuyingRate)
	}
	if resp.Data.SellingRate != 130.75 {
		t.Errorf("SellingRate = %v, want 130.75", resp.Data.SellingRate)
	}
}

func TestClient_DecodesAPIErrorOn4xx(t *testing.T) {
	client, ts := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"code":422,"message":"Transaction hash has already been processed"}`))
	})
	defer ts.Close()

	_, err := client.ExchangeRate(context.Background(), ExchangeRateBody{CurrencyCode: KES})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err is not *APIError: %T (%v)", err, err)
	}
	if apiErr.Code != 422 {
		t.Errorf("Code = %d, want 422", apiErr.Code)
	}
	if !errors.Is(err, ErrTransactionHashAlreadyProcessed) {
		t.Errorf("errors.Is(err, ErrTransactionHashAlreadyProcessed) = false, want true")
	}
}

func TestClient_FallsBackToStatusCodeWhenAPIErrorCodeIsZero(t *testing.T) {
	client, ts := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"bad input"}`))
	})
	defer ts.Close()

	_, err := client.ExchangeRate(context.Background(), ExchangeRateBody{CurrencyCode: KES})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err is not *APIError: %T (%v)", err, err)
	}
	if apiErr.Code != http.StatusBadRequest {
		t.Errorf("Code = %d, want %d (fallback to HTTP status)", apiErr.Code, http.StatusBadRequest)
	}
	if apiErr.Message != "bad input" {
		t.Errorf("Message = %q, want bad input", apiErr.Message)
	}
}

func TestClient_NonJSONErrorSurfacesStatus(t *testing.T) {
	client, ts := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("upstream offline"))
	})
	defer ts.Close()

	_, err := client.ExchangeRate(context.Background(), ExchangeRateBody{CurrencyCode: KES})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("err = %q, want it to include the HTTP status", err.Error())
	}
	// NB: parseResponse consumes the body via json.Decode before falling back
	// to io.ReadAll, so the upstream body text is currently lost for non-JSON
	// error responses. Asserting only the status keeps this test stable if
	// that bug gets fixed.
}
