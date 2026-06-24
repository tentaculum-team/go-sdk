package payments

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrNotFound is returned when the requested resource does not exist. The
	// high-level wrappers (e.g. ActiveSubscription) translate it into a nil
	// result where "absent" is a normal outcome.
	ErrNotFound = errors.New("payments: not found")
	// ErrUnauthorized is returned on a 401 from payments.
	ErrUnauthorized = errors.New("payments: unauthorized")
	// ErrForbidden is returned on a 403 from payments.
	ErrForbidden = errors.New("payments: forbidden")
)

// APIError wraps any non-2xx the SDK did not map to a sentinel. The server's
// envelope message is preserved in Message.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("payments: api error (status %d)", e.StatusCode)
	}
	return fmt.Sprintf("payments: api error (status %d): %s", e.StatusCode, e.Message)
}

// mapAPIError translates an HTTP status + envelope message into a sentinel,
// falling back to *APIError.
func mapAPIError(status int, message string) error {
	switch status {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	}
	return &APIError{StatusCode: status, Message: message}
}
