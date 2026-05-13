package pretium

import (
	"strings"
	"testing"
)

func TestParseWebhook_StatusVariant(t *testing.T) {
	body := `{
		"status": "COMPLETE",
		"transaction_code": "e37c02ca-2170-4a82-ad06-d2def781cc8e",
		"receipt_number": "TKTQRBEO7A",
		"public_name": "John Doe",
		"message": "Transaction processed successfully."
	}`

	got, err := ParseWebhook(strings.NewReader(body))
	if err != nil {
		t.Fatalf("ParseWebhook returned error: %v", err)
	}

	if got.Event() != WebhookEventStatus {
		t.Errorf("Event() = %q, want %q", got.Event(), WebhookEventStatus)
	}
	if got.Status != StatusComplete {
		t.Errorf("Status = %q, want %q", got.Status, StatusComplete)
	}
	if got.TransactionCode != "e37c02ca-2170-4a82-ad06-d2def781cc8e" {
		t.Errorf("TransactionCode = %q", got.TransactionCode)
	}
	if got.ReceiptNumber == nil || *got.ReceiptNumber != "TKTQRBEO7A" {
		t.Errorf("ReceiptNumber = %v, want TKTQRBEO7A", got.ReceiptNumber)
	}
	if got.PublicName == nil || *got.PublicName != "John Doe" {
		t.Errorf("PublicName = %v, want John Doe", got.PublicName)
	}
	if got.IsReleased != nil {
		t.Errorf("IsReleased = %v, want nil for status variant", got.IsReleased)
	}
	if got.TransactionHash != nil {
		t.Errorf("TransactionHash = %v, want nil for status variant", got.TransactionHash)
	}
}

func TestParseWebhook_AssetReleasedVariant(t *testing.T) {
	body := `{
		"is_released": true,
		"transaction_code": "e37c02ca-2170-4a82-ad06-d2def781cc8e",
		"transaction_hash": "0xabc123"
	}`

	got, err := ParseWebhook(strings.NewReader(body))
	if err != nil {
		t.Fatalf("ParseWebhook returned error: %v", err)
	}

	if got.Event() != WebhookEventAssetReleased {
		t.Errorf("Event() = %q, want %q", got.Event(), WebhookEventAssetReleased)
	}
	if got.IsReleased == nil || !*got.IsReleased {
		t.Errorf("IsReleased = %v, want pointer to true", got.IsReleased)
	}
	if got.TransactionHash == nil || *got.TransactionHash != "0xabc123" {
		t.Errorf("TransactionHash = %v, want 0xabc123", got.TransactionHash)
	}
	if got.ReceiptNumber != nil {
		t.Errorf("ReceiptNumber = %v, want nil for release variant", got.ReceiptNumber)
	}
	if got.Status != "" {
		t.Errorf("Status = %q, want empty for release variant", got.Status)
	}
}

func TestParseWebhook_IsReleasedFalseStillCountsAsReleaseEvent(t *testing.T) {
	body := `{"is_released": false, "transaction_code": "tc", "transaction_hash": "0xdead"}`

	got, err := ParseWebhook(strings.NewReader(body))
	if err != nil {
		t.Fatalf("ParseWebhook returned error: %v", err)
	}
	if got.Event() != WebhookEventAssetReleased {
		t.Errorf("Event() = %q, want %q (is_released present discriminates by presence, not value)", got.Event(), WebhookEventAssetReleased)
	}
}

func TestParseWebhook_Malformed(t *testing.T) {
	_, err := ParseWebhook(strings.NewReader("{not json"))
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}
