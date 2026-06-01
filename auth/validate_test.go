package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const testSecret = "test-access-secret"

// signAccess builds an HS256 access token shaped like auth-api Claims.
func signAccess(t *testing.T, secret string, ttl time.Duration, method gojwt.SigningMethod) string {
	t.Helper()
	cl := gojwt.MapClaims{
		"user_id":   uuid.New().String(),
		"org_id":    uuid.New().String(),
		"user_type": "user",
		"is_owner":  true,
		"role":      "ADMIN",
		"exp":       time.Now().Add(ttl).Unix(),
		"iat":       time.Now().Unix(),
		"jti":       uuid.New().String(),
	}
	tok := gojwt.NewWithClaims(method, cl)
	var key any = []byte(secret)
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return s
}

func TestValidateTokenLocal(t *testing.T) {
	c, _ := New(Config{AccessSecret: testSecret})

	t.Run("valid", func(t *testing.T) {
		tok := signAccess(t, testSecret, 15*time.Minute, gojwt.SigningMethodHS256)
		id, err := c.ValidateTokenLocal(tok)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if id.UserType != UserTypePersonal || id.Role != RoleAdmin || !id.IsOwner {
			t.Fatalf("bad identity: %+v", id)
		}
		if id.ExpiresAt.IsZero() {
			t.Fatal("ExpiresAt should be set from exp")
		}
	})

	t.Run("expired", func(t *testing.T) {
		tok := signAccess(t, testSecret, -time.Minute, gojwt.SigningMethodHS256)
		_, err := c.ValidateTokenLocal(tok)
		if !errors.Is(err, ErrTokenExpired) {
			t.Fatalf("want ErrTokenExpired, got %v", err)
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		tok := signAccess(t, "other-secret", 15*time.Minute, gojwt.SigningMethodHS256)
		_, err := c.ValidateTokenLocal(tok)
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
	uid, oid := uuid.New(), uuid.New()
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
			`"user_id":"` + uid.String() + `","org_id":"` + oid.String() + `",` +
			`"user_type":"enterprise_user","is_owner":false,"role":"USER",` +
			`"email":"a@b.com","username":"alice"}}`))
	}))
	defer srv.Close()

	c, _ := New(Config{BaseURL: srv.URL})

	t.Run("ok", func(t *testing.T) {
		id, err := c.ValidateToken(context.Background(), "good")
		if err != nil {
			t.Fatalf("unexpected: %v", err)
		}
		if id.UserID != uid || id.Email != "a@b.com" || id.Username != "alice" {
			t.Fatalf("bad identity: %+v", id)
		}
		if id.UserType != UserTypeEnterprise {
			t.Fatalf("bad user_type: %s", id.UserType)
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
		"totp_required":            ErrTOTPRequired,
		"account_pending_deletion": ErrAccountPendingDeletion,
		"oauth_account":            ErrOAuthAccount,
		"invalid credentials":      ErrInvalidCredentials,
		"invalid token":            ErrInvalidToken,
		"missing refresh token":    ErrMissingRefreshToken,
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
