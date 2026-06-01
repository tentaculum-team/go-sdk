package auth

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidToken           = errors.New("auth: invalid token")
	ErrTokenExpired           = errors.New("auth: token expired")
	ErrTOTPRequired           = errors.New("auth: totp required")
	ErrAccountPendingDeletion = errors.New("auth: account pending deletion")
	ErrOAuthAccount           = errors.New("auth: oauth account, use provider login")
	ErrOAuthEmailUnverified   = errors.New("auth: oauth email not verified by provider")
	ErrOAuthLinkRequired      = errors.New("auth: email already registered, link provider from an authenticated session")
	ErrOAuthAlreadyLinked     = errors.New("auth: provider already linked to another account")
	ErrInvalidCredentials     = errors.New("auth: invalid credentials")
	ErrMissingRefreshToken    = errors.New("auth: missing refresh token")
	ErrOfflineDisabled        = errors.New("auth: offline validation disabled (no AccessSecret)")
	ErrInternalDisabled       = errors.New("auth: service tokens disabled (no InternalSecret)")
	ErrNoBaseURL              = errors.New("auth: BaseURL required for remote calls")
)

// APIError wraps any non-2xx the SDK didn't map to a sentinel.
type APIError struct {
	StatusCode int
	Message    string // server's `message` field
}

func (e *APIError) Error() string {
	return fmt.Sprintf("auth: api error (status %d): %s", e.StatusCode, e.Message)
}

// mapAPIError translates the server envelope message (machine codes from
// auth-api) into a sentinel, falling back to *APIError.
func mapAPIError(status int, message string) error {
	switch message {
	case "totp_required":
		return ErrTOTPRequired
	case "account_pending_deletion":
		return ErrAccountPendingDeletion
	case "oauth_account":
		return ErrOAuthAccount
	case "oauth_email_unverified":
		return ErrOAuthEmailUnverified
	case "oauth_link_required":
		return ErrOAuthLinkRequired
	case "oauth_already_linked":
		return ErrOAuthAlreadyLinked
	case "invalid credentials":
		return ErrInvalidCredentials
	case "invalid token", "missing bearer token":
		return ErrInvalidToken
	case "missing refresh token":
		return ErrMissingRefreshToken
	}
	return &APIError{StatusCode: status, Message: message}
}
