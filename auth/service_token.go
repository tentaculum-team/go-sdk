package auth

import (
	"errors"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// serviceTokenTTL matches auth-api jwt.GenerateServiceToken (30s).
const serviceTokenTTL = 30 * time.Second

// ServiceClaims is the verified payload of a service token (registered
// claims only — no identity).
type ServiceClaims struct {
	ID        string // jti
	ExpiresAt time.Time
	IssuedAt  time.Time
}

// GenerateServiceToken mints a short-lived (30s) HS256 service token signed
// with InternalSecret, registered-claims only with a random jti. Mirrors
// auth-api jwt.GenerateServiceToken.
//
// NOTE: the service does not yet verify these (gap §11.2); until it adopts
// VerifyServiceToken/RequireServiceToken, service-to-service calls are
// unauthenticated.
func (c *Client) GenerateServiceToken() (string, error) {
	if c.internalSecret == nil {
		return "", ErrInternalDisabled
	}
	now := time.Now()
	tok := gojwt.NewWithClaims(gojwt.SigningMethodHS256, gojwt.RegisteredClaims{
		ExpiresAt: gojwt.NewNumericDate(now.Add(serviceTokenTTL)),
		IssuedAt:  gojwt.NewNumericDate(now),
		ID:        uuid.New().String(),
	})
	return tok.SignedString(c.internalSecret)
}

// VerifyServiceToken checks HS256 signature (InternalSecret) and exp, and
// returns the claims (incl. jti for an optional replay guard). The 30s TTL
// is the skew budget between caller and callee clocks.
func (c *Client) VerifyServiceToken(raw string) (*ServiceClaims, error) {
	if c.internalSecret == nil {
		return nil, ErrInternalDisabled
	}
	var rc gojwt.RegisteredClaims
	token, err := gojwt.ParseWithClaims(raw, &rc, func(t *gojwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return c.internalSecret, nil
	})
	if err != nil {
		if errors.Is(err, gojwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	sc := &ServiceClaims{ID: rc.ID}
	if rc.ExpiresAt != nil {
		sc.ExpiresAt = rc.ExpiresAt.Time
	}
	if rc.IssuedAt != nil {
		sc.IssuedAt = rc.IssuedAt.Time
	}
	return sc, nil
}
