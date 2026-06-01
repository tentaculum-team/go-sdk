// Package httpmw provides framework-agnostic net/http middleware mirroring
// auth-api WithAuth / OwnerOnly / AdminOnly, backed by the auth SDK Client.
//
// Import path: github.com/Tentaculum-dev/go-sdk/auth/middleware/nethttp
package httpmw

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Tentaculum-dev/go-sdk/auth"
	"github.com/google/uuid"
)

type ctxKeyType struct{}

var identityKey = ctxKeyType{}

type config struct {
	headerTrust bool
	local       bool
}

// Option configures WithAuth.
type Option func(*config)

// WithHeaderTrust enables the unsigned X-* gateway header path. Default OFF;
// only safe behind a trusted gateway. (§5.3)
func WithHeaderTrust() Option { return func(c *config) { c.headerTrust = true } }

// WithLocalValidation validates tokens offline (requires AccessSecret).
func WithLocalValidation() Option { return func(c *config) { c.local = true } }

// WithAuth wraps next, authenticating the request and storing an
// *auth.Identity in the request context. On failure it writes 401 with the
// service envelope shape.
func WithAuth(c *auth.Client, next http.Handler, opts ...Option) http.Handler {
	var cfg config
	for _, o := range opts {
		o(&cfg)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfg.headerTrust {
			if id, ok := identityFromHeaders(r); ok {
				serveWith(next, w, r, id)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			write401(w)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		var (
			id  *auth.Identity
			err error
		)
		if cfg.local {
			id, err = c.ValidateTokenLocal(token)
		} else {
			id, err = c.ValidateToken(r.Context(), token)
		}
		if err != nil {
			write401(w)
			return
		}
		serveWith(next, w, r, id)
	})
}

// OwnerOnly wraps next, requiring the stored Identity to be owner.
func OwnerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := IdentityFromContext(r.Context())
		if !ok || !id.IsOwner {
			write403(w, "owner only")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AdminOnly wraps next, requiring the stored Identity to be ADMIN.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := IdentityFromContext(r.Context())
		if !ok || id.Role != auth.RoleAdmin {
			write403(w, "admin only")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireRole wraps next, requiring the stored Identity to hold one of roles.
func RequireRole(next http.Handler, roles ...auth.Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := IdentityFromContext(r.Context())
		if !ok {
			write403(w, "forbidden")
			return
		}
		for _, role := range roles {
			if id.Role == role {
				next.ServeHTTP(w, r)
				return
			}
		}
		write403(w, "forbidden")
	})
}

// RequireServiceToken verifies an internal service token from X-Service-Token.
func RequireServiceToken(c *auth.Client, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get("X-Service-Token")
		if raw == "" {
			write401(w)
			return
		}
		if _, err := c.VerifyServiceToken(raw); err != nil {
			write401(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// IdentityFromContext returns the Identity stored by WithAuth.
func IdentityFromContext(ctx context.Context) (*auth.Identity, bool) {
	id, ok := ctx.Value(identityKey).(*auth.Identity)
	return id, ok
}

func serveWith(next http.Handler, w http.ResponseWriter, r *http.Request, id *auth.Identity) {
	ctx := context.WithValue(r.Context(), identityKey, id)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func identityFromHeaders(r *http.Request) (*auth.Identity, bool) {
	userID, errU := uuid.Parse(r.Header.Get("X-User-ID"))
	orgID, errO := uuid.Parse(r.Header.Get("X-Org-ID"))
	if errU != nil || errO != nil {
		return nil, false
	}
	return &auth.Identity{
		UserID:   userID,
		OrgID:    orgID,
		UserType: auth.UserType(r.Header.Get("X-User-Type")),
		IsOwner:  r.Header.Get("X-Is-Owner") == "true",
		Role:     auth.Role(r.Header.Get("X-User-Role")),
	}, true
}

func write401(w http.ResponseWriter) { writeEnvelope(w, http.StatusUnauthorized, "unauthorized") }

func write403(w http.ResponseWriter, msg string) { writeEnvelope(w, http.StatusForbidden, msg) }

func writeEnvelope(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "message": msg})
}
