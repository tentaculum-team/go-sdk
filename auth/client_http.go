package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// ── inputs (mirror auth-api DTOs; binding tags are server-side) ───

type RegisterInput struct {
	OrgName   string `json:"org_name"`
	OrgSlug   string `json:"org_slug"`
	Plan      string `json:"plan"` // "user" | "enterprise" (NOT the user_type claim space)
	SeatLimit *int   `json:"seat_limit,omitempty"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TOTPCode string `json:"totp_code,omitempty"`
}

// ── outputs ───

// AuthResult is the `data` of login/confirm/refresh (AuthOutput).
type AuthResult struct {
	UserID      string  `json:"user_id"`
	OrgID       string  `json:"org_id"`
	UserType    string  `json:"user_type"`
	IsOwner     bool    `json:"is_owner"`
	Role        string  `json:"role"`
	AccessToken string  `json:"access_token"`
	AvatarUUID  *string `json:"avatar_uuid,omitempty"`
}

// User mirrors auth-api UserOutput (GET /users/me).
type User struct {
	ID            string  `json:"id"`
	OrgID         string  `json:"org_id"`
	Email         string  `json:"email"`
	Username      string  `json:"username"`
	FirstName     *string `json:"first_name,omitempty"`
	LastName      *string `json:"last_name,omitempty"`
	Phone         *string `json:"phone,omitempty"`
	Bio           *string `json:"bio,omitempty"`
	AvatarUUID    *string `json:"avatar_uuid,omitempty"`
	Department    string  `json:"department,omitempty"`
	JobTitle      *string `json:"job_title,omitempty"`
	EmployeeID    *string `json:"employee_id,omitempty"`
	Locale        string  `json:"locale,omitempty"`
	Timezone      string  `json:"timezone,omitempty"`
	BirthDate     *string `json:"birth_date,omitempty"`
	Country       *string `json:"country,omitempty"`
	HasPassword   bool    `json:"has_password"`
	TOTPEnabled   bool    `json:"totp_enabled"`
	IsActive      bool    `json:"is_active"`
	EmailVerified bool    `json:"email_verified"`
	LastLoginAt   *string `json:"last_login_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// Organization mirrors auth-api OrgOutput (GET /organizations/me).
type Organization struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Plan        string  `json:"plan"`
	OwnerID     string  `json:"owner_id"`
	LogoUUID    *string `json:"logo_uuid,omitempty"`
	Description *string `json:"description,omitempty"`
	Website     *string `json:"website,omitempty"`
	Country     *string `json:"country,omitempty"`
	Industry    *string `json:"industry,omitempty"`
	SeatLimit   *int    `json:"seat_limit,omitempty"`
	SeatUsed    int     `json:"seat_used"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ── endpoint wrappers

// Login authenticates a user. May return ErrTOTPRequired (422) — re-call
// with TOTPCode set. Returns the rotated refresh_token cookie for forwarding.
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
// the refresh_token COOKIE (not the body), so it is sent as a cookie.
// Returns the new AuthResult and the rotated refresh_token cookie.
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

// Org fetches the authenticated org (GET /organizations/me) with a Bearer token.
func (c *Client) Org(ctx context.Context, accessToken string) (*Organization, error) {
	req, err := c.bearerGET(ctx, c.apiBase+"/organizations/me", accessToken)
	if err != nil {
		return nil, err
	}
	env, err := c.doJSON(req, true)
	if err != nil {
		return nil, err
	}
	var out Organization
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ── internals ─

func (c *Client) bearerGET(ctx context.Context, url, accessToken string) (*http.Request, error) {
	req, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return req, nil
}

// doJSONCookies posts a JSON body (when non-nil), attaches request cookies,
// and returns the envelope plus the response Set-Cookie list.
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
