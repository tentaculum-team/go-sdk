package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// validateOutput is the `data` shape of GET /auth/validate.
type validateOutput struct {
	UserID     string  `json:"user_id"`
	OrgID      string  `json:"org_id"`
	UserType   string  `json:"user_type"`
	IsOwner    bool    `json:"is_owner"`
	Role       string  `json:"role"`
	Email      string  `json:"email"`
	Username   string  `json:"username"`
	AvatarUUID *string `json:"avatar_uuid,omitempty"`
}

// ValidateToken validates an access token remotely via GET /auth/validate.
// No secret sharing; returns email/username; honors central revocation.
//
// Maps 401 -> ErrInvalidToken. Retries an idempotent GET once on a
// connection error (never on 401). Uses the configured cache if set.
func (c *Client) ValidateToken(ctx context.Context, accessToken string) (*Identity, error) {
	if c.baseURL == "" {
		return nil, ErrNoBaseURL
	}

	var key string
	if c.cache != nil {
		key = cacheKey(accessToken)
		if id, ok := c.cache.Get(key); ok {
			return id, nil
		}
	}

	req, err := c.newRequest(ctx, http.MethodGet, c.apiBase+"/auth/validate", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	env, err := c.doJSON(req, true)
	if err != nil {
		return nil, err
	}

	var out validateOutput
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, err
	}
	id, err := out.toIdentity()
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		c.cache.Set(key, id, c.cacheCap)
	}
	c.log.Debug("token validated remotely", "user_id", id.UserID, "org_id", id.OrgID)
	return id, nil
}

// ValidateTokenLocal validates an access token offline using AccessSecret
// (== service JWT_SECRET). Zero network. Cannot see server-side revocation;
// a logged-out access token stays valid until exp (<=15m), same as the
// service's own middleware.
//
// Returns ErrOfflineDisabled if AccessSecret is empty, ErrTokenExpired on
// expiry, ErrInvalidToken otherwise. Email/Username are not populated.
func (c *Client) ValidateTokenLocal(accessToken string) (*Identity, error) {
	if c.accessSecret == nil {
		return nil, ErrOfflineDisabled
	}
	return parseAccess(accessToken, c.accessSecret)
}

// parseAccess mirrors auth-api pkg/jwt.parse: HS256 only, reject non-HMAC.
func parseAccess(raw string, secret []byte) (*Identity, error) {
	var cl claims
	token, err := gojwt.ParseWithClaims(raw, &cl, func(t *gojwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
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
	return cl.toIdentity(), nil
}

func (o validateOutput) toIdentity() (*Identity, error) {
	uid, err := uuid.Parse(o.UserID)
	if err != nil {
		return nil, err
	}
	oid, err := uuid.Parse(o.OrgID)
	if err != nil {
		return nil, err
	}
	id := &Identity{
		UserID:   uid,
		OrgID:    oid,
		UserType: UserType(o.UserType),
		IsOwner:  o.IsOwner,
		Role:     Role(o.Role),
		Email:    o.Email,
		Username: o.Username,
	}
	if o.AvatarUUID != nil && *o.AvatarUUID != "" {
		if av, err := uuid.Parse(*o.AvatarUUID); err == nil {
			id.AvatarUUID = &av
		}
	}
	return id, nil
}

func cacheKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
