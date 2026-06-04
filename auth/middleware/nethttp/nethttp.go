// Package httpmw provides framework-agnostic net/http middleware mirroring
// auth-api authentication, backed by the auth SDK Client.
//
// Import path: github.com/Tentaculum-dev/go-sdk/auth/middleware/nethttp
package httpmw

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Tentaculum-dev/go-sdk/auth"
)

type ctxKeyType struct{}

var identityKey = ctxKeyType{}

type config struct {
	local bool
}

// Option configures WithAuth.
type Option func(*config)

// WithLocalValidation validates tokens offline (requires AccessSecret).
func WithLocalValidation() Option { return func(c *config) { c.local = true } }

// WithAuth wraps next, authenticating the request (Bearer token) and storing an
// *auth.Identity in the request context. On failure it writes 401 with the
// service envelope shape.
func WithAuth(c *auth.Client, next http.Handler, opts ...Option) http.Handler {
	var cfg config
	for _, o := range opts {
		o(&cfg)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

// RequireRole wraps next, requiring the stored Identity's sys_role to be one of
// roles (e.g. auth.RoleAdmin).
func RequireRole(next http.Handler, roles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := IdentityFromContext(r.Context())
		if !ok {
			write403(w, "forbidden")
			return
		}
		for _, role := range roles {
			if id.SysRole == role {
				next.ServeHTTP(w, r)
				return
			}
		}
		write403(w, "forbidden")
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

func write401(w http.ResponseWriter) { writeEnvelope(w, http.StatusUnauthorized, "unauthorized") }

func write403(w http.ResponseWriter, msg string) { writeEnvelope(w, http.StatusForbidden, msg) }

func writeEnvelope(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "message": msg})
}
