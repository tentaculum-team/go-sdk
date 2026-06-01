# auth SDK

Go SDK for the `tentaculum-auth` service (`../auth`). Validate end-user tokens,
protect routes, read the authenticated identity, and do service-to-service auth.

Import path: `github.com/ViitoJooj/sdk/auth`

## Install

```go
import (
    "github.com/ViitoJooj/sdk/auth"
    ginmw  "github.com/ViitoJooj/sdk/auth/middleware/gin"
    httpmw "github.com/ViitoJooj/sdk/auth/middleware/nethttp"
    "github.com/ViitoJooj/sdk/auth/cache"
)
```

Gin is only pulled in if you import `middleware/gin`. Non-Gin services use
`middleware/nethttp` and never compile Gin.

## Config & environments (dev/prod)

`auth.ConfigFromEnv()` is the only place the SDK reads env:

| Env var               | Maps to          | Notes |
|-----------------------|------------------|-------|
| `AUTH_ENV`            | env selection    | `prod`/`production` → prod, else dev |
| `AUTH_URL_PROD`       | `BaseURL` (prod) | used when `AUTH_ENV=prod` |
| `AUTH_URL_DEV`        | `BaseURL` (dev)  | used otherwise |
| `AUTH_URL`            | `BaseURL`        | fallback if per-env unset |
| `JWT_SECRET`          | `AccessSecret`   | enables offline validation |
| `INTERNAL_JWT_SECRET` | `InternalSecret` | enables service tokens |
| `AUTH_USER_AGENT`     | `UserAgent`      | |

```go
cfg := auth.ConfigFromEnv()
cfg.Cache = cache.NewLRU(1024) // optional remote-validation cache
client, _ := auth.New(cfg)
```

Or construct directly:

```go
client, _ := auth.New(auth.Config{
    BaseURL:      "https://auth.internal:8080",
    AccessSecret: os.Getenv("JWT_SECRET"), // optional: offline path
})
```

`Client` is safe for concurrent use.

## Token validation

**Remote (default).** No secret sharing; returns email/username; honors central
revocation. One network hop (cache to soften).

```go
id, err := client.ValidateToken(ctx, accessToken)
// errors.Is(err, auth.ErrInvalidToken)
```

**Offline.** Requires `AccessSecret == service JWT_SECRET`. Zero network. Cannot
see revocation — a logged-out access token stays valid until `exp` (≤15m), same
as the service's own middleware. No email/username.

```go
id, err := client.ValidateTokenLocal(accessToken)
// auth.ErrTokenExpired / auth.ErrInvalidToken / auth.ErrOfflineDisabled
```

Use remote for sensitive flows; offline for high-RPS internal calls where a
≤15-min revocation window is acceptable.

## Middleware

### Gin

```go
api := r.Group("/api/v1", ginmw.WithAuth(client))
api.GET("/things", func(c *gin.Context) {
    id, _ := ginmw.IdentityFrom(c)
    _ = id.UserID
})
owner := api.Group("/admin", ginmw.OwnerOnly())
admin := api.Group("/ops", ginmw.AdminOnly())
api.GET("/r", ginmw.RequireRole(auth.RoleAdmin), h)
```

Options: `WithLocalValidation()`, `WithHeaderTrust()`, `WithContextKey(k)`.

### net/http

```go
h := httpmw.WithAuth(client, next)            // store Identity in ctx
h = httpmw.OwnerOnly(h)
id, ok := httpmw.IdentityFromContext(r.Context())
```

Failure responses match the service envelope:
`401 {"success":false,"message":"unauthorized"}`,
`403 {"success":false,"message":"owner only"|"admin only"|"forbidden"}`.

### Gateway header mode (`X-*`) — opt-in, unsigned ⚠️

`WithHeaderTrust()` replicates the service's behavior: when **both** `X-User-ID`
and `X-Org-ID` are valid UUIDs, identity is trusted from headers **without a
signature**; otherwise it falls through to token validation. Only enable behind a
trusted gateway that strips/sets these headers on every request. **Default OFF.**

## Service-to-service

```go
tok, _ := client.GenerateServiceToken()         // 30s, INTERNAL_JWT_SECRET
_, err := client.VerifyServiceToken(tok)
r.GET("/internal", ginmw.RequireServiceToken(client), h) // header: X-Service-Token
```

> The service does not yet verify these (gap). Until it adopts
> `VerifyServiceToken`/`RequireServiceToken`, service-to-service calls are
> unauthenticated.

## HTTP wrappers (proxying)

```go
res, refreshCookie, err := client.Login(ctx, auth.LoginInput{Email, Password})
// err may be auth.ErrTOTPRequired -> re-call with TOTPCode
err = client.Register(ctx, auth.RegisterInput{...})        // 202, email sent
res, refreshCookie, err = client.Refresh(ctx, refreshToken) // sent as cookie
err = client.Logout(ctx, refreshToken)                      // sent as cookie
user, err := client.Me(ctx, accessToken)
org, err := client.Org(ctx, accessToken)
```

`Refresh`/`Logout` send the refresh token as the `refresh_token` **cookie** (the
service reads it there, not the body) and return the rotated `Set-Cookie` to
forward.

## Errors

Sentinels for `errors.Is`: `ErrInvalidToken`, `ErrTokenExpired`,
`ErrTOTPRequired`, `ErrAccountPendingDeletion`, `ErrOAuthAccount`,
`ErrInvalidCredentials`, `ErrMissingRefreshToken`, `ErrOfflineDisabled`,
`ErrInternalDisabled`, `ErrNoBaseURL`. Anything unmapped → `*auth.APIError`
(`StatusCode`, `Message`).

## Notes on the service

- `user_type` claim space (`user`|`enterprise_user`) differs from the register
  `plan` field (`user`|`enterprise`). The SDK exposes claim values verbatim.
- `GET /auth/validate` is Bearer-only (ignores cookie and `X-*`).
- No `/api-keys` routes exist (table dropped); the SDK exposes none.
