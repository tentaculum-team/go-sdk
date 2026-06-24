// Package driver is a Go SDK for the tentaculum driver (file storage) service.
// It is used server-side to store and remove files. The driver trusts the
// X-User-UUID / X-Sys-Role headers (the same identity model the protect gateway
// injects), so this client sets them explicitly — it must only ever reach the
// driver over the internal network, never be exposed to end users directly.
package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// Driver storage buckets (mirror the driver's domain).
const (
	BucketUser    = "user"
	BucketCompany = "company"
	BucketSite    = "site"
	BucketPublic  = "public"
)

const (
	// Version is the SDK version.
	Version = "0.1.0"

	defaultTimeout = 15 * time.Second
)

// Client talks to the driver service. Safe for concurrent use.
type Client struct {
	apiBase   string // baseURL + "/api/v1"
	http      *http.Client
	userAgent string
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient overrides the default *http.Client (15s timeout).
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

// New builds a Client. baseURL is the driver root, e.g. "http://driver:8082".
func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		apiBase:   strings.TrimRight(baseURL, "/") + "/api/v1",
		http:      &http.Client{Timeout: defaultTimeout},
		userAgent: "driver-sdk-go/" + Version,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// FileRef is the subset of a driver file callers care about.
type FileRef struct {
	Uuid string
	URL  string
}

// UploadInput describes a file to store in the driver.
type UploadInput struct {
	OwnerUUID   string // becomes the file owner (X-User-UUID + ownership)
	SysRole     string // X-Sys-Role (defaults to "USER")
	Bucket      string // user | company | site | public (defaults to "user")
	Filename    string
	ContentType string // optional; multipart part content-type
	IsPublic    bool
	Data        []byte
}

// envelope mirrors the driver success/error response shape.
type envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Uuid        string `json:"uuid"`
		DownloadURL string `json:"download_url"`
	} `json:"data"`
}

// Upload stores a file in the driver and returns its reference. Multipart POST
// to /api/v1/files with the owner injected as X-User-UUID.
func (c *Client) Upload(ctx context.Context, in UploadInput) (*FileRef, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	part, err := w.CreateFormFile("file", in.Filename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(in.Data); err != nil {
		return nil, err
	}
	bucket := in.Bucket
	if bucket == "" {
		bucket = BucketUser
	}
	_ = w.WriteField("bucket", bucket)
	if in.IsPublic {
		_ = w.WriteField("is_public", "true")
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiBase+"/files", &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("User-Agent", c.userAgent)
	c.setIdentity(req, in.OwnerUUID, in.SysRole)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var out envelope
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil, fmt.Errorf("driver upload: decode response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, mapAPIError(resp.StatusCode, out.Message)
	}
	return &FileRef{Uuid: out.Data.Uuid, URL: out.Data.DownloadURL}, nil
}

// Delete removes a file owned by ownerUUID. DELETE /api/v1/files/:uuid.
func (c *Client) Delete(ctx context.Context, ownerUUID, sysRole, fileUUID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.apiBase+"/files/"+fileUUID, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", c.userAgent)
	c.setIdentity(req, ownerUUID, sysRole)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var env envelope
		_ = json.NewDecoder(resp.Body).Decode(&env)
		return mapAPIError(resp.StatusCode, env.Message)
	}
	return nil
}

func (c *Client) setIdentity(req *http.Request, ownerUUID, sysRole string) {
	req.Header.Set("X-User-UUID", ownerUUID)
	if sysRole == "" {
		sysRole = "USER"
	}
	req.Header.Set("X-Sys-Role", sysRole)
}
