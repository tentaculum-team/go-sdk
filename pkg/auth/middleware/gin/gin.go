// Package ginmw provides Gin middleware mirroring the auth-api authentication,
// backed by the auth SDK Client.
//
// Import path: github.com/Tentaculum-dev/go-sdk/auth/middleware/gin
package ginmw

import (
	"net/http"
	"strings"

	"github.com/Tentaculum-dev/go-sdk/pkg/auth"
	"github.com/gin-gonic/gin"
)

const contextKey = "auth.identity"

type config struct {
	local bool
}

// Option configures WithAuth.
type Option func(*config)

// WithLocalValidation validates tokens offline via Client.ValidateTokenLocal
// instead of the remote endpoint. Requires AccessSecret on the Client.
func WithLocalValidation() Option { return func(c *config) { c.local = true } }

// WithAuth authenticates the request (Bearer token) and stores an
// *auth.Identity in the gin context. On failure it aborts with 401 and the
// service envelope shape.
func WithAuth(c *auth.Client, opts ...Option) gin.HandlerFunc {
	var cfg config
	for _, o := range opts {
		o(&cfg)
	}

	return func(gc *gin.Context) {
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

		gc.Set(contextKey, id)
		gc.Next()
	}
}

// RequireRole aborts with 403 unless the stored Identity's sys_role is one of
// roles (e.g. auth.RoleAdmin).
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(gc *gin.Context) {
		id, ok := IdentityFrom(gc)
		if !ok {
			abort403(gc, "forbidden")
			return
		}
		for _, r := range roles {
			if id.SysRole == r {
				gc.Next()
				return
			}
		}
		abort403(gc, "forbidden")
	}
}

// IdentityFrom returns the Identity stored by WithAuth.
func IdentityFrom(gc *gin.Context) (*auth.Identity, bool) {
	v, ok := gc.Get(contextKey)
	if !ok {
		return nil, false
	}
	id, ok := v.(*auth.Identity)
	return id, ok
}

func abort401(gc *gin.Context) {
	gc.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
}

func abort403(gc *gin.Context, msg string) {
	gc.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "message": msg})
}
