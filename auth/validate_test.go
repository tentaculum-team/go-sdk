package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
)

func signAccess(secret paseto.V4AsymmetricSecretKey, ttl time.Duration) string {
	t := paseto.NewToken()
	now := time.Now()
	t.SetIssuedAt(now)
	t.SetExpiration(now.Add(ttl))
	t.SetString("typ", "access")
	t.SetString("user_uuid", "u1")
	t.SetString("sys_role", "ADMIN")
	return t.V4Sign(secret, nil)
}

func TestValidateTokenLocal(t *testing.T) {
	secret := paseto.NewV4AsymmetricSecretKey()
	c, _ := New(Config{AccessPublicKey: secret.Public().ExportHex()})

	t.Run("valid", func(t *testing.T) {
		id, err := c.ValidateTokenLocal(signAccess(secret, 15*time.Minute))
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if id.UserUUID != "u1" || id.SysRole != RoleAdmin {
			t.Fatalf("bad identity: %+v", id)
		}
		if id.ExpiresAt.IsZero() {
			t.Fatal("ExpiresAt should be set")
		}
	})

	t.Run("expired", func(t *testing.T) {
		_, err := c.ValidateTokenLocal(signAccess(secret, -time.Minute))
		if !errors.Is(err, ErrTokenExpired) {
			t.Fatalf("want ErrTokenExpired, got %v", err)
		}
	})

	t.Run("wrong key", func(t *testing.T) {
		other := paseto.NewV4AsymmetricSecretKey()
		_, err := c.ValidateTokenLocal(signAccess(other, 15*time.Minute))
		if !errors.Is(err, ErrInvalidToken) {
			t.Fatalf("want ErrInvalidToken, got %v", err)
		}
	})

	t.Run("offline disabled", func(t *testing.T) {
		nc, _ := New(Config{})
		_, err := nc.ValidateTokenLocal("x")
		if !errors.Is(err, ErrOfflineDisabled) {
			t.Fatalf("want ErrOfflineDisabled, got %v", err)
		}
	})
}

func TestValidateTokenRemote(t *testing.T) {
	uid := "0190a6c2-0000-7000-8000-000000000000"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/auth/validate" {
			t.Errorf("bad path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer good" {
			w.WriteHeader(401)
			_, _ = w.Write([]byte(`{"success":false,"message":"invalid token"}`))
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"success":true,"message":"valid","data":{` +
			`"user_uuid":"` + uid + `","sys_role":"USER",` +
			`"email":"a@b.com","username":"alice"}}`))
	}))
	defer srv.Close()

	c, _ := New(Config{BaseURL: srv.URL})

	t.Run("ok", func(t *testing.T) {
		id, err := c.ValidateToken(context.Background(), "good")
		if err != nil {
			t.Fatalf("unexpected: %v", err)
		}
		if id.UserUUID != uid || id.Email != "a@b.com" || id.Username != "alice" {
			t.Fatalf("bad identity: %+v", id)
		}
		if id.SysRole != RoleUser {
			t.Fatalf("bad sys_role: %s", id.SysRole)
		}
	})

	t.Run("invalid -> ErrInvalidToken", func(t *testing.T) {
		_, err := c.ValidateToken(context.Background(), "bad")
		if !errors.Is(err, ErrInvalidToken) {
			t.Fatalf("want ErrInvalidToken, got %v", err)
		}
	})

	t.Run("no base url", func(t *testing.T) {
		nc, _ := New(Config{})
		_, err := nc.ValidateToken(context.Background(), "x")
		if !errors.Is(err, ErrNoBaseURL) {
			t.Fatalf("want ErrNoBaseURL, got %v", err)
		}
	})
}

func TestMapAPIError(t *testing.T) {
	cases := map[string]error{
		"oauth_account":         ErrOAuthAccount,
		"oauth_link_required":   ErrOAuthLinkRequired,
		"invalid credentials":   ErrInvalidCredentials,
		"invalid token":         ErrInvalidToken,
		"missing refresh token": ErrMissingRefreshToken,
	}
	for msg, want := range cases {
		if got := mapAPIError(400, msg); !errors.Is(got, want) {
			t.Errorf("%q -> %v, want %v", msg, got, want)
		}
	}
	var apiErr *APIError
	if got := mapAPIError(500, "boom"); !errors.As(got, &apiErr) || apiErr.StatusCode != 500 {
		t.Errorf("unknown message should map to *APIError, got %v", got)
	}
}
