package pretium

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PretiumClient struct {
	settlementAddress string
	apiKey            string
	httpClient        *http.Client
}

const (
	userAgent   = "pretium-go"
	contentType = "application/json"

	baseLiveEndpoint = "https://api.xwift.africa"
	versionPath      = "/v1/"
)

// New returns an instance of a Pretium client reusable across different products
func New(settlementAddress string, apiKey string) *PretiumClient {
	PretiumClient := &PretiumClient{
		settlementAddress: settlementAddress,
		apiKey:            apiKey,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	return PretiumClient
}

// SetHTTPClient can be used to override the default client with a custom set one
func (fc *PretiumClient) SetHTTPClient(httpClient *http.Client) *PretiumClient {
	fc.httpClient = httpClient

	return fc
}

// setHeaders sets the headers required by the Fonbnk API
func (fc *PretiumClient) setHeaders(req *http.Request) (*http.Request, error) {
	if err := fc.setAuthHeaders(req); err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)

	return req, nil
}

// requestWithCtx builds the HTTP request
func (fc *PretiumClient) requestWithCtx(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	return fc.do(req)
}

// do executes the built http request, setting appropriate headers
func (fc *PretiumClient) do(req *http.Request) (*http.Response, error) {
	builtRequest, err := fc.setHeaders(req)
	if err != nil {
		return nil, err
	}

	return fc.httpClient.Do(builtRequest)
}

// parseResponse is a general utility to decode JSON responses correctly
func parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("Pretium server error: code=%s: response_body=%s", resp.Status, string(b))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
