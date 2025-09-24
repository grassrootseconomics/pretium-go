package pretium

import (
	"errors"
	"fmt"
	"strings"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

var (
	ErrTransactionHashAlreadyProcessed = errors.New("transaction hash already processed")
	ErrAmountMismatch                  = errors.New("amount mismatch equivalent amount in USD")
	ErrAmountBelowMinimum              = errors.New("amount below minimum")
)

// Pretium doesn't support error codes, so we match based on error message content
func (e *APIError) Is(target error) bool {
	switch target {
	case ErrTransactionHashAlreadyProcessed:
		return strings.Contains(strings.ToLower(e.Message), "already been processed")
	case ErrAmountMismatch:
		m := strings.ToLower(e.Message)
		return strings.Contains(m, "mismatch equivalent amount in usd") || strings.Contains(m, "mismatch equivalent amount")
	case ErrAmountBelowMinimum:
		return strings.HasPrefix(strings.ToLower(e.Message), "the amount field must be at least")
	default:
		return false
	}
}
