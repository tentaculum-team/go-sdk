package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"aidanwoods.dev/go-paseto"
)

// validateOutput is the `data` shape of GET /auth/validate.
type validateOutput struct {
	UserUUID string  `json:"user_uuid"`
	SysRole  string  `json:"sys_role"`
	Email    string  `json:"email"`
	Username string  `json:"username"`
	ImgURL   *string `json:"img_url,omitempty"`
}

// ValidateToken validates an access token remotely via GET /auth/validate.
// No key sharing; returns email/username; honors central revocation.
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
	id := out.toIdentity()

	if c.cache != nil {
		c.cache.Set(key, id, c.cacheCap)
	}
	c.log.Debug("token validated remotely", "user_uuid", id.UserUUID)
	return id, nil
}

// ValidateTokenLocal validates an access token offline using the auth service's
// PASETO v4.public key (AccessPublicKey). Zero network. Cannot see server-side
// revocation; a logged-out access token stays valid until exp.
//
// Returns ErrOfflineDisabled if AccessPublicKey is empty, ErrTokenExpired on
// expiry, ErrInvalidToken otherwise. Email/Username are not populated.
func (c *Client) ValidateTokenLocal(accessToken string) (*Identity, error) {
	if c.accessPublicKey == "" {
		return nil, ErrOfflineDisabled
	}
	return parseAccess(accessToken, c.accessPublicKey)
}

// parseAccess verifies a PASETO v4.public access token with the public key.
func parseAccess(raw, pubHex string) (*Identity, error) {
	key, err := paseto.NewV4AsymmetricPublicKeyFromHex(pubHex)
	if err != nil {
		return nil, ErrInvalidToken
	}
	parser := paseto.NewParser() // NotExpired rule preloaded
	t, err := parser.ParseV4Public(key, raw, nil)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}
	if typ, _ := t.GetString("typ"); typ != "access" {
		return nil, ErrInvalidToken
	}

	id := &Identity{}
	id.UserUUID, _ = t.GetString("user_uuid")
	id.SysRole, _ = t.GetString("sys_role")
	if exp, err := t.GetExpiration(); err == nil {
		id.ExpiresAt = exp
	}
	return id, nil
}

func (o validateOutput) toIdentity() *Identity {
	return &Identity{
		UserUUID: o.UserUUID,
		SysRole:  o.SysRole,
		Email:    o.Email,
		Username: o.Username,
		ImgURL:   o.ImgURL,
	}
}

func cacheKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
