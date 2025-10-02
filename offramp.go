package pretium

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type (
	ExchangeRateBody struct {
		CurrencyCode string `json:"currency_code"`
	}

	ExchangeRateResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			BuyingRate  string `json:"buying_rate"`
			SellingRate string `json:"selling_rate"`
		} `json:"data"`
	}

	ValidationBody struct {
		Type          string `json:"type"`
		Shortcode     string `json:"shortcode"`
		MobileNetwork string `json:"mobile_network"`
	}

	ValidationResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Status        string `json:"status"`
			Shortcode     string `json:"shortcode"`
			PublicName    string `json:"public_name"`
			MobileNetwork string `json:"mobile_network"`
		} `json:"data"`
	}

	PayBody struct {
		TransactionHash string `json:"transaction_hash"`
		Amount          string `json:"amount"`
		Shortcode       string `json:"shortcode"`
		Type            string `json:"type"`
		Chain           string `json:"chain"`
	}

	PayResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Status          string `json:"status"`
			TransactionCode string `json:"transaction_code"`
			Message         string `json:"message"`
		} `json:"data"`
	}
	StatusResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			ID              int     `json:"id"`
			TransactionCode string  `json:"transaction_code"`
			Status          string  `json:"status"`
			Amount          string  `json:"amount"`
			AmountInUSD     string  `json:"amount_in_usd"`
			Type            string  `json:"type"`
			Shortcode       string  `json:"shortcode"`
			AccountNumber   *string `json:"account_number"`
			PublicName      *string `json:"public_name"`
			ReceiptNumber   *string `json:"receipt_number"`
			Category        string  `json:"category"`
			Chain           string  `json:"chain"`
			Asset           string  `json:"asset"`
			TransactionHash *string `json:"transaction_hash"`
			Message         string  `json:"message"`
			CurrencyCode    string  `json:"currency_code"`
			IsReleased      bool    `json:"is_released"`
			CreatedAt       string  `json:"created_at"`
		} `json:"data"`
	}

	StatusBody struct {
		TransactionCode string `json:"transaction_code"`
	}

	WebhookPayload struct {
		Status          string `json:"status"`
		TransactionCode string `json:"transaction_code"`
		Message         string `json:"message"`
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

func (fc *PretiumClient) Validation(ctx context.Context, input ValidationBody) (ValidationResponse, error) {
	validationResp := ValidationResponse{}

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

func (fc *PretiumClient) Pay(ctx context.Context, input PayBody) (PayResponse, error) {
	payResp := PayResponse{}

	payload := struct {
		PayBody
		CallbackURL string `json:"callback_url"`
	}{
		PayBody:     input,
		CallbackURL: fc.callbackURL,
	}

	b, err := json.Marshal(&payload)
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

func (fc *PretiumClient) Status(ctx context.Context, input StatusBody) (StatusResponse, error) {
	statusResp := StatusResponse{}

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

func ParseWebhook(r io.Reader) (WebhookPayload, error) {
	var webhook WebhookPayload
	if err := json.NewDecoder(r).Decode(&webhook); err != nil {
		return webhook, err
	}
	return webhook, nil
}
