package pretium

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type (
	OnrampBody struct {
		Shortcode     string  `json:"shortcode"`
		Amount        float64 `json:"amount"`
		MobileNetwork string  `json:"mobile_network"`
		Chain         string  `json:"chain"`
		Asset         string  `json:"asset"`
		Address       string  `json:"address"`
	}

	OnrampResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			TransactionCode string `json:"transaction_code"`
			Status          string `json:"status"`
			Message         string `json:"message"`
		} `json:"data"`
	}
)

func (fc *PretiumClient) Onramp(ctx context.Context, currencyCode string, input OnrampBody) (OnrampResponse, error) {
	onrampResp := OnrampResponse{}

	payload := struct {
		OnrampBody
		CallbackURL string `json:"callback_url"`
	}{
		OnrampBody:  input,
		CallbackURL: fc.callbackURL,
	}

	b, err := json.Marshal(&payload)
	if err != nil {
		return onrampResp, err
	}

	url := baseLiveEndpoint + versionPath + "onramp/" + currencyCode
	resp, err := fc.requestWithCtx(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return onrampResp, err
	}

	if err := parseResponse(resp, &onrampResp); err != nil {
		return onrampResp, err
	}

	return onrampResp, nil
}
