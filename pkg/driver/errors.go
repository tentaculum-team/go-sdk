package driver

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrNoBaseURL is returned when the client was built without a base URL.
	ErrNoBaseURL = errors.New("driver: baseURL required")
	// ErrNotFound is returned when the target file does not exist.
	ErrNotFound = errors.New("driver: file not found")
	// ErrUnauthorized is returned on a 401 from the driver.
	ErrUnauthorized = errors.New("driver: unauthorized")
	// ErrForbidden is returned on a 403 from the driver.
	ErrForbidden = errors.New("driver: forbidden")
)

// APIError wraps any non-2xx the SDK did not map to a sentinel. The server's
// envelope message is preserved in Message.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("driver: api error (status %d)", e.StatusCode)
	}
	return fmt.Sprintf("driver: api error (status %d): %s", e.StatusCode, e.Message)
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
