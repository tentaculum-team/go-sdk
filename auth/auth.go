// Package auth is a Go SDK for the tentaculum-auth service. It validates
// end-user access tokens (remote or offline), exposes the authenticated
// identity, handles service-to-service tokens, and provides thin HTTP
// wrappers for the public auth endpoints.
//
// Middleware lives in subpackages so non-Gin consumers don't pull Gin:
//   - github.com/Tentaculum-dev/go-sdk/auth/middleware/gin
//   - github.com/Tentaculum-dev/go-sdk/auth/middleware/nethttp
package auth

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	// Version is the SDK version, appended to the default User-Agent.
	Version = "0.1.0"

	defaultTimeout = 5 * time.Second
	defaultAPIVer  = "v1"
)

// TokenCache is an opt-in cache for remote validation results, keyed by a
// hash of the token. Implementations must be safe for concurrent use.
type TokenCache interface {
	Get(key string) (*Identity, bool)
	Set(key string, id *Identity, ttl time.Duration)
}

// Config configures a Client. The zero value is not usable; call New.
type Config struct {
	// BaseURL of the auth service, e.g. "https://auth.internal:8080".
	// Required for remote validation and the HTTP client wrappers.
	BaseURL string

	// APIVersion pins the base path segment. Defaults to "v1".
	APIVersion string

	// HTTPClient is optional; defaults to a client with a 5s timeout.
	HTTPClient *http.Client

	// AccessSecret enables OFFLINE token validation (== service JWT_SECRET).
	// If empty, ValidateTokenLocal returns ErrOfflineDisabled.
	AccessSecret string

	// InternalSecret enables service-to-service token gen/verify
	// (== service INTERNAL_JWT_SECRET). Optional.
	InternalSecret string

	// UserAgent appended to outbound requests. Defaults to "auth-sdk-go/<ver>".
	UserAgent string

	// Cache, if set, memoizes remote ValidateToken results. Default: none.
	// TTL = min(remaining exp, CacheCap).
	Cache TokenCache

	// CacheCap caps cache entry TTL. Defaults to 60s when Cache is set.
	CacheCap time.Duration

	// Logger logs at debug level only. Tokens/secrets are never logged.
	Logger *slog.Logger
}

// Client talks to the auth service. Safe for concurrent use.
type Client struct {
	baseURL        string
	apiBase        string // baseURL + "/api/" + version
	http           *http.Client
	accessSecret   []byte
	internalSecret []byte
	userAgent      string
	cache          TokenCache
	cacheCap       time.Duration
	log            *slog.Logger
}

// New builds a Client. BaseURL is required unless the client is used purely
// for offline validation / service tokens.
func New(cfg Config) (*Client, error) {
	httpc := cfg.HTTPClient
	if httpc == nil {
		httpc = &http.Client{Timeout: defaultTimeout}
	}
	ver := cfg.APIVersion
	if ver == "" {
		ver = defaultAPIVer
	}
	ua := cfg.UserAgent
	if ua == "" {
		ua = "auth-sdk-go/" + Version
	}
	base := strings.TrimRight(cfg.BaseURL, "/")
	cacheCap := cfg.CacheCap
	if cacheCap <= 0 {
		cacheCap = 60 * time.Second
	}
	log := cfg.Logger
	if log == nil {
		log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	c := &Client{
		baseURL:   base,
		apiBase:   base + "/api/" + ver,
		http:      httpc,
		userAgent: ua,
		cache:     cfg.Cache,
		cacheCap:  cacheCap,
		log:       log,
	}
	if cfg.AccessSecret != "" {
		c.accessSecret = []byte(cfg.AccessSecret)
	}
	if cfg.InternalSecret != "" {
		c.internalSecret = []byte(cfg.InternalSecret)
	}
	return c, nil
}

// envelope is the shared response shape: OKResponse / ErrorResponse.
type envelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// doJSON performs an HTTP request and decodes the envelope. On non-2xx it
// returns a mapped error (sentinel or *APIError). retryGET retries once on
// connection-level errors for idempotent GETs.
func (c *Client) doJSON(req *http.Request, retryGET bool) (*envelope, error) {
	if c.baseURL == "" {
		return nil, ErrNoBaseURL
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		if retryGET && req.Method == http.MethodGet && req.GetBody != nil {
			if body, berr := req.GetBody(); berr == nil {
				req.Body = body
				resp, err = c.http.Do(req)
			}
		} else if retryGET && req.Method == http.MethodGet {
			resp, err = c.http.Do(req)
		}
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()

	var env envelope
	raw, _ := io.ReadAll(resp.Body)
	if len(raw) > 0 {
		if jerr := json.Unmarshal(raw, &env); jerr != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil, jerr
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, mapAPIError(resp.StatusCode, env.Message)
	}
	return &env, nil
}

func (c *Client) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

// HasOffline reports whether offline validation is configured.
func (c *Client) HasOffline() bool { return c.accessSecret != nil }
