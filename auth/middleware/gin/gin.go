// Package ginmw provides Gin middleware mirroring the auth-api WithAuth /
// OwnerOnly / AdminOnly behavior, backed by the auth SDK Client.
//
// Import path: github.com/ViitoJooj/sdk/auth/middleware/gin
package ginmw

import (
	"net/http"
	"strings"

	"github.com/ViitoJooj/sdk/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ctxKeyType is the private default storage key for the Identity.
type ctxKeyType struct{}

var defaultKey = ctxKeyType{}

type config struct {
	headerTrust bool
	local       bool
	key         any
}

// Option configures WithAuth.
type Option func(*config)

// WithHeaderTrust enables the X-* gateway header path. MUST only be used when
// the service sits behind a trusted gateway that strips/sets these headers on
// every request — they are unsigned. Default OFF. (§5.3)
func WithHeaderTrust() Option { return func(c *config) { c.headerTrust = true } }

// WithLocalValidation validates tokens offline via Client.ValidateTokenLocal
// instead of the remote endpoint. Requires AccessSecret on the Client.
func WithLocalValidation() Option { return func(c *config) { c.local = true } }

// WithContextKey overrides the storage key. NOTE: the guard helpers
// (OwnerOnly/AdminOnly/RequireRole) and IdentityFrom read the default key;
// if you override it, read the Identity yourself.
func WithContextKey(k any) Option { return func(c *config) { c.key = k } }

// WithAuth authenticates the request and stores an *auth.Identity in the gin
// context. On failure it aborts with 401 and the service envelope shape.
func WithAuth(c *auth.Client, opts ...Option) gin.HandlerFunc {
	cfg := config{key: defaultKey}
	for _, o := range opts {
		o(&cfg)
	}

	return func(gc *gin.Context) {
		// 1. Gateway header trust (opt-in): both UUIDs valid -> trust.
		if cfg.headerTrust {
			if id, ok := identityFromHeaders(gc); ok {
				gc.Set(keyString(cfg.key), id)
				gc.Next()
				return
			}
		}

		// 2. Bearer token.
		authHeader := gc.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			abort401(gc)
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
			id, err = c.ValidateToken(gc.Request.Context(), token)
		}
		if err != nil {
			abort401(gc)
			return
		}

		gc.Set(keyString(cfg.key), id)
		gc.Next()
	}
}

// OwnerOnly aborts with 403 "owner only" unless the stored Identity is owner.
func OwnerOnly() gin.HandlerFunc {
	return func(gc *gin.Context) {
		id, ok := IdentityFrom(gc)
		if !ok || !id.IsOwner {
			abort403(gc, "owner only")
			return
		}
		gc.Next()
	}
}

// AdminOnly aborts with 403 "admin only" unless the stored Identity is ADMIN.
func AdminOnly() gin.HandlerFunc {
	return func(gc *gin.Context) {
		id, ok := IdentityFrom(gc)
		if !ok || id.Role != auth.RoleAdmin {
			abort403(gc, "admin only")
			return
		}
		gc.Next()
	}
}

// RequireRole aborts with 403 unless the stored Identity holds one of roles.
func RequireRole(roles ...auth.Role) gin.HandlerFunc {
	return func(gc *gin.Context) {
		id, ok := IdentityFrom(gc)
		if !ok {
			abort403(gc, "forbidden")
			return
		}
		for _, r := range roles {
			if id.Role == r {
				gc.Next()
				return
			}
		}
		abort403(gc, "forbidden")
	}
}

// IdentityFrom returns the Identity stored under the default key.
func IdentityFrom(gc *gin.Context) (*auth.Identity, bool) {
	v, ok := gc.Get(keyString(defaultKey))
	if !ok {
		return nil, false
	}
	id, ok := v.(*auth.Identity)
	return id, ok
}

// identityFromHeaders builds an Identity from X-* gateway headers. Returns
// ok=false unless both X-User-ID and X-Org-ID parse as UUIDs (mirrors
// auth-api middleware.WithAuth).
func identityFromHeaders(gc *gin.Context) (*auth.Identity, bool) {
	userID, errU := uuid.Parse(gc.GetHeader("X-User-ID"))
	orgID, errO := uuid.Parse(gc.GetHeader("X-Org-ID"))
	if errU != nil || errO != nil {
		return nil, false
	}
	return &auth.Identity{
		UserID:   userID,
		OrgID:    orgID,
		UserType: auth.UserType(gc.GetHeader("X-User-Type")),
		IsOwner:  gc.GetHeader("X-Is-Owner") == "true",
		Role:     auth.Role(gc.GetHeader("X-User-Role")),
	}, true
}

// keyString turns any key into the string form gin.Context.Set expects.
// The default key is a stable private constant.
func keyString(k any) string {
	if k == defaultKey {
		return "auth.identity"
	}
	if s, ok := k.(string); ok {
		return s
	}
	// Fallback: derive a deterministic-ish string. Custom non-string keys
	// are discouraged for gin (Set takes a string).
	return "auth.identity"
}

// RequireServiceToken verifies an internal service token from the
// X-Service-Token header (kept separate from Authorization to avoid colliding
// with end-user bearer tokens). Aborts 401 on missing/invalid token.
// Requires InternalSecret on the Client.
func RequireServiceToken(c *auth.Client) gin.HandlerFunc {
	return func(gc *gin.Context) {
		raw := gc.GetHeader("X-Service-Token")
		if raw == "" {
			abort401(gc)
			return
		}
		if _, err := c.VerifyServiceToken(raw); err != nil {
			abort401(gc)
			return
		}
		gc.Next()
	}
}

func abort401(gc *gin.Context) {
	gc.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
}

func abort403(gc *gin.Context, msg string) {
	gc.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "message": msg})
}
