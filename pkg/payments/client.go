// Package payments is a Go SDK for the tentaculum payments service. It exposes
// the product catalog and the billing gate other services need (is the caller's
// company subscribed?). The payments service trusts the X-User-UUID / X-Sys-Role
// headers the protect gateway injects, so this client sets them explicitly and
// must only reach payments over the internal network.
package payments

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// Version is the SDK version.
	Version = "0.1.0"

	defaultTimeout = 10 * time.Second
)

// Client talks to the payments service. Safe for concurrent use.
type Client struct {
	apiBase   string // baseURL + "/api/v1"
	http      *http.Client
	userAgent string
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient overrides the default *http.Client (10s timeout).
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.http = h
		}
	}
}

// WithUserAgent overrides the User-Agent sent on requests.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		if ua != "" {
			c.userAgent = ua
		}
	}
}

// New builds a Client. baseURL is the payments root, e.g. "http://payments:18082".
func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		apiBase:   strings.TrimRight(baseURL, "/") + "/api/v1",
		http:      &http.Client{Timeout: defaultTimeout},
		userAgent: "payments-sdk-go/" + Version,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Product mirrors a payments catalog entry. Price is in cents.
type Product struct {
	Uuid            string  `json:"uuid"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Type            string  `json:"type"` // "software" | "plan"
	Price           int64   `json:"price"`
	Currency        string  `json:"currency"`
	BillingInterval *string `json:"billing_interval,omitempty"`
	Active          bool    `json:"active"`
	ImageURL        *string `json:"image_url,omitempty"`
}

// Subscription mirrors a payments subscription (company-scoped).
type Subscription struct {
	Uuid               string    `json:"uuid"`
	CompanyUuid        string    `json:"company_uuid"`
	PlanUuid           string    `json:"plan_uuid"`
	Provider           string    `json:"provider"`
	Status             string    `json:"status"`
	CurrentPeriodEnd   time.Time `json:"current_period_end"`
	CancelAtPeriodEnd  bool      `json:"cancel_at_period_end"`
}

// envelope is the shared payments response shape.
type envelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// ListProducts returns the catalog. typ filters by "software"/"plan" (empty =
// all); active filters by active flag (nil = all).
func (c *Client) ListProducts(ctx context.Context, typ string, active *bool) ([]Product, error) {
	q := url.Values{}
	if typ != "" {
		q.Set("type", typ)
	}
	if active != nil {
		q.Set("active", strconv.FormatBool(*active))
	}
	u := c.apiBase + "/products"
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	env, err := c.get(ctx, u, "", "")
	if err != nil {
		return nil, err
	}
	var out []Product
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetProduct fetches a single product by UUID.
func (c *Client) GetProduct(ctx context.Context, uuid string) (*Product, error) {
	env, err := c.get(ctx, c.apiBase+"/products/"+uuid, "", "")
	if err != nil {
		return nil, err
	}
	var out Product
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ActiveSubscription returns the caller's active subscription, or nil if none.
// Payments resolves the user's company internally from the X-User-UUID header.
// GET /api/v1/subscriptions/active — 404 means "no active subscription".
func (c *Client) ActiveSubscription(ctx context.Context, userUUID, sysRole string) (*Subscription, error) {
	env, err := c.get(ctx, c.apiBase+"/subscriptions/active", userUUID, sysRole)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	var out Subscription
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) get(ctx context.Context, u, userUUID, sysRole string) (*envelope, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if userUUID != "" {
		req.Header.Set("X-User-UUID", userUUID)
		if sysRole == "" {
			sysRole = "USER"
		}
		req.Header.Set("X-Sys-Role", sysRole)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var env envelope
	_ = json.NewDecoder(resp.Body).Decode(&env)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, mapAPIError(resp.StatusCode, env.Message)
	}
	return &env, nil
}
