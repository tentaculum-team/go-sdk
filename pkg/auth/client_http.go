package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// ── inputs (mirror auth-api request bodies) ───

// RegisterInput is the body of POST /auth/register. A company is a separate,
// optional resource created later — registration only creates a personal user.
type RegisterInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewRegister builds a registration body.
func NewRegister(email, username, password string) RegisterInput {
	return RegisterInput{Email: email, Username: username, Password: password}
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ── outputs ───

// AuthResult is the `data` of login/confirm/refresh.
type AuthResult struct {
	UserUUID    string  `json:"user_uuid"`
	SysRole     string  `json:"sys_role"`
	AccessToken string  `json:"access_token"`
	ImgURL      *string `json:"img_url,omitempty"`
}

// User mirrors auth-api GET /users/me.
type User struct {
	Uuid            string  `json:"uuid"`
	ImgURL          *string `json:"img_url,omitempty"`
	Username        string  `json:"username"`
	Email           string  `json:"email"`
	Phone           *string `json:"phone,omitempty"`
	SysRole         string  `json:"sys_role"`
	FirstName       *string `json:"first_name,omitempty"`
	LastName        *string `json:"last_name,omitempty"`
	EmailVerifiedAt *string `json:"email_verified_at,omitempty"`
	GoogleAuth      bool    `json:"google_auth"`
	GitHubAuth      bool    `json:"github_auth"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// ── endpoint wrappers ───

// Login authenticates a user and returns the rotated refresh_token cookie for
// forwarding.
func (c *Client) Login(ctx context.Context, in LoginInput) (*AuthResult, *http.Cookie, error) {
	env, cookies, err := c.doJSONCookies(ctx, http.MethodPost, c.apiBase+"/auth/login", in, nil)
	if err != nil {
		return nil, nil, err
	}
	var out AuthResult
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, nil, err
	}
	return &out, findCookie(cookies, "refresh_token"), nil
}

// Register starts a pending registration (202, confirmation email sent).
func (c *Client) Register(ctx context.Context, in RegisterInput) error {
	_, _, err := c.doJSONCookies(ctx, http.MethodPost, c.apiBase+"/auth/register", in, nil)
	return err
}

// Refresh rotates the access token. The service reads the refresh token from
// the refresh_token COOKIE (not the body). Returns the new AuthResult and the
// rotated refresh_token cookie.
func (c *Client) Refresh(ctx context.Context, refreshToken string) (*AuthResult, *http.Cookie, error) {
	jar := []*http.Cookie{{Name: "refresh_token", Value: refreshToken}}
	env, cookies, err := c.doJSONCookies(ctx, http.MethodPost, c.apiBase+"/auth/token", nil, jar)
	if err != nil {
		return nil, nil, err
	}
	var out AuthResult
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, nil, err
	}
	return &out, findCookie(cookies, "refresh_token"), nil
}

// Logout revokes the session. The refresh token is sent as a cookie.
func (c *Client) Logout(ctx context.Context, refreshToken string) error {
	jar := []*http.Cookie{{Name: "refresh_token", Value: refreshToken}}
	_, _, err := c.doJSONCookies(ctx, http.MethodPost, c.apiBase+"/auth/logout", nil, jar)
	return err
}

// Me fetches the authenticated user (GET /users/me) with a Bearer token.
func (c *Client) Me(ctx context.Context, accessToken string) (*User, error) {
	req, err := c.bearerGET(ctx, c.apiBase+"/users/me", accessToken)
	if err != nil {
		return nil, err
	}
	env, err := c.doJSON(req, true)
	if err != nil {
		return nil, err
	}
	var out User
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ── OAuth (browser redirect URLs) ───

// OAuth providers supported by the auth service.
const (
	ProviderGoogle = "google"
	ProviderGitHub = "github"
)

// OAuthLoginURL builds the URL to start an OAuth login for the given provider
// ("google"/"github"). Redirect the user's browser here; the service sends them
// back to redirectURI with the session (or ?erro=<code> on failure).
func (c *Client) OAuthLoginURL(provider, redirectURI string) string {
	return c.oauthURL(provider, "", redirectURI)
}

// OAuthLinkURL builds the URL to link a provider to the already-logged-in user.
// The browser must carry the access_token cookie (set by the auth service); the
// service reads it to identify who to link. On success it returns to
// redirectURI with ?oauth=linked&provider=<provider>.
func (c *Client) OAuthLinkURL(provider, redirectURI string) string {
	return c.oauthURL(provider, "link", redirectURI)
}

func (c *Client) oauthURL(provider, action, redirectURI string) string {
	path := c.apiBase + "/auth/oauth/" + provider
	if action != "" {
		path += "/" + action
	}
	if redirectURI != "" {
		path += "?redirect_uri=" + url.QueryEscape(redirectURI)
	}
	return path
}

// ── internals ───

func (c *Client) bearerGET(ctx context.Context, url, accessToken string) (*http.Request, error) {
	req, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return req, nil
}

// doJSONCookies posts a JSON body (when non-nil), attaches request cookies, and
// returns the envelope plus the response Set-Cookie list.
func (c *Client) doJSONCookies(ctx context.Context, method, url string, body any, cookies []*http.Cookie) (*envelope, []*http.Cookie, error) {
	if c.baseURL == "" {
		return nil, nil, ErrNoBaseURL
	}
	var rdr *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		rdr = bytes.NewReader(b)
	} else {
		rdr = bytes.NewReader(nil)
	}
	req, err := c.newRequest(ctx, method, url, rdr)
	if err != nil {
		return nil, nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	respCookies := resp.Cookies()
	var env envelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil, respCookies, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, respCookies, mapAPIError(resp.StatusCode, env.Message)
	}
	return &env, respCookies, nil
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, ck := range cookies {
		if ck.Name == name {
			return ck
		}
	}
	return nil
}
