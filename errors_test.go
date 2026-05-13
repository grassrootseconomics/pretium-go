package pretium

import (
	"errors"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	e := &APIError{Code: 422, Message: "boom"}
	if got, want := e.Error(), "422: boom"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestAPIError_Is(t *testing.T) {
	tests := []struct {
		name    string
		message string
		target  error
		want    bool
	}{
		{
			name:    "transaction hash already processed - exact phrasing",
			message: "Transaction hash has already been processed",
			target:  ErrTransactionHashAlreadyProcessed,
			want:    true,
		},
		{
			name:    "transaction hash already processed - mixed case",
			message: "this hash HAS ALREADY BEEN PROCESSED please retry",
			target:  ErrTransactionHashAlreadyProcessed,
			want:    true,
		},
		{
			name:    "amount mismatch - usd phrasing",
			message: "Amount mismatch equivalent amount in USD",
			target:  ErrAmountMismatch,
			want:    true,
		},
		{
			name:    "amount mismatch - generic phrasing",
			message: "Amount mismatch equivalent amount",
			target:  ErrAmountMismatch,
			want:    true,
		},
		{
			name:    "amount below minimum",
			message: "The amount field must be at least 10",
			target:  ErrAmountBelowMinimum,
			want:    true,
		},
		{
			name:    "amount below minimum - lowercase prefix",
			message: "the amount field must be at least 1000",
			target:  ErrAmountBelowMinimum,
			want:    true,
		},
		{
			name:    "unrelated message",
			message: "internal server error",
			target:  ErrAmountBelowMinimum,
			want:    false,
		},
		{
			name:    "amount below minimum prefix is required - substring is not enough",
			message: "validation failed: the amount field must be at least 10",
			target:  ErrAmountBelowMinimum,
			want:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := &APIError{Code: 400, Message: tc.message}
			if got := errors.Is(e, tc.target); got != tc.want {
				t.Errorf("errors.Is(%q, %v) = %v, want %v", tc.message, tc.target, got, tc.want)
			}
		})
	}
}

func TestAPIError_Is_UnknownTargetReturnsFalse(t *testing.T) {
	e := &APIError{Code: 400, Message: "anything"}
	other := errors.New("some other sentinel")
	if errors.Is(e, other) {
		t.Error("errors.Is should return false for unrelated target")
	}
}
