package pretium

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type (
	ExchangeRateBody struct {
		CurrencyCode string `json:"currency_code"`
	}

	ExchangeRateResponse struct {
		Success bool `json:"success"`
	}

	ValidationBody struct {
		Type          string `json:"type"`
		Shortcode     string `json:"shortcode"`
		MobileNetwork string `json:"mobile_network"`
	}

	PayBody struct {
		TransactionHash string `json:"transaction_hash"`
		Amount          string `json:"amount"`
		Shortcode       string `json:"shortcode"`
		Type            string `json:"type"`
		Chain           string `json:"chain"`
	}

	StatusBody struct {
		TransactionCode string `json:"transaction_code"`
	}
)

func (fc *PretiumClient) ExchangeRate(ctx context.Context, input ExchangeRateBody) (ExchangeRateResponse, error) {
	exchangeRateResp := ExchangeRateResponse{}

	jsonRequestBody, err := json.Marshal(&input)
	if err != nil {
		return exchangeRateResp, err
	}

	resp, err := fc.requestWithCtx(ctx, http.MethodPost, baseLiveEndpoint+versionPath+"exchange-rate", bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		return exchangeRateResp, err
	}

	if err := parseResponse(resp, &exchangeRateResp); err != nil {
		return exchangeRateResp, err
	}

	return exchangeRateResp, nil
}

func (fc *PretiumClient) Validation(ctx context.Context, input ValidationBody) (ExchangeRateResponse, error) {
	validationResp := ExchangeRateResponse{}

	b, err := json.Marshal(&input)
	if err != nil {
		return validationResp, err
	}

	resp, err := fc.requestWithCtx(ctx, http.MethodPost, baseLiveEndpoint+versionPath+"validation", bytes.NewBuffer(b))
	if err != nil {
		return validationResp, err
	}

	if err := parseResponse(resp, &validationResp); err != nil {
		return validationResp, err
	}

	return validationResp, nil
}

func (fc *PretiumClient) Pay(ctx context.Context, input PayBody) (ExchangeRateResponse, error) {
	payResp := ExchangeRateResponse{}

	b, err := json.Marshal(&input)
	if err != nil {
		return payResp, err
	}

	resp, err := fc.requestWithCtx(ctx, http.MethodPost, baseLiveEndpoint+versionPath+"pay", bytes.NewBuffer(b))
	if err != nil {
		return payResp, err
	}

	if err := parseResponse(resp, &payResp); err != nil {
		return payResp, err
	}

	return payResp, nil
}

func (fc *PretiumClient) Status(ctx context.Context, input StatusBody) (ExchangeRateResponse, error) {
	statusResp := ExchangeRateResponse{}

	b, err := json.Marshal(&input)
	if err != nil {
		return statusResp, err
	}

	resp, err := fc.requestWithCtx(ctx, http.MethodPost, baseLiveEndpoint+versionPath+"status", bytes.NewBuffer(b))
	if err != nil {
		return statusResp, err
	}

	if err := parseResponse(resp, &statusResp); err != nil {
		return statusResp, err
	}

	return statusResp, nil
}
