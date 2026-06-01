package ginmw

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Tentaculum-dev/go-sdk/auth"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const secret = "s3cr3t"

func init() { gin.SetMode(gin.TestMode) }

func sign(t *testing.T, role string, owner bool) string {
	t.Helper()
	cl := gojwt.MapClaims{
		"user_id": uuid.New().String(), "org_id": uuid.New().String(),
		"user_type": "user", "is_owner": owner, "role": role,
		"exp": time.Now().Add(time.Minute).Unix(), "iat": time.Now().Unix(),
	}
	s, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, cl).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func newClient(t *testing.T) *auth.Client {
	c, err := auth.New(auth.Config{AccessSecret: secret})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestWithAuthLocal(t *testing.T) {
	c := newClient(t)
	r := gin.New()
	r.GET("/x", WithAuth(c, WithLocalValidation()), func(gc *gin.Context) {
		id, ok := IdentityFrom(gc)
		if !ok {
			gc.JSON(500, gin.H{})
			return
		}
		gc.JSON(200, gin.H{"role": string(id.Role)})
	})

	t.Run("valid bearer", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/x", nil)
		req.Header.Set("Authorization", "Bearer "+sign(t, "ADMIN", true))
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatalf("code %d body %s", w.Code, w.Body.String())
		}
	})

	t.Run("missing bearer -> 401 envelope", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/x", nil)
		r.ServeHTTP(w, req)
		if w.Code != 401 {
			t.Fatalf("want 401 got %d", w.Code)
		}
		var env map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &env)
		if env["success"] != false || env["message"] != "unauthorized" {
			t.Fatalf("bad envelope: %s", w.Body.String())
		}
	})
}

func TestHeaderTrust(t *testing.T) {
	c := newClient(t)
	uid, oid := uuid.New(), uuid.New()

	t.Run("disabled ignores headers -> 401", func(t *testing.T) {
		r := gin.New()
		r.GET("/x", WithAuth(c, WithLocalValidation()), func(gc *gin.Context) { gc.JSON(200, gin.H{}) })
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/x", nil)
		req.Header.Set("X-User-ID", uid.String())
		req.Header.Set("X-Org-ID", oid.String())
		r.ServeHTTP(w, req)
		if w.Code != 401 {
			t.Fatalf("want 401 got %d", w.Code)
		}
	})

	t.Run("enabled trusts valid UUIDs", func(t *testing.T) {
		r := gin.New()
		r.GET("/x", WithAuth(c, WithHeaderTrust()), func(gc *gin.Context) {
			id, _ := IdentityFrom(gc)
			gc.JSON(200, gin.H{"u": id.UserID.String()})
		})
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/x", nil)
		req.Header.Set("X-User-ID", uid.String())
		req.Header.Set("X-Org-ID", oid.String())
		req.Header.Set("X-User-Role", "ADMIN")
		req.Header.Set("X-Is-Owner", "true")
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatalf("want 200 got %d body %s", w.Code, w.Body.String())
		}
	})

	t.Run("enabled but malformed falls through to bearer", func(t *testing.T) {
		r := gin.New()
		r.GET("/x", WithAuth(c, WithHeaderTrust(), WithLocalValidation()), func(gc *gin.Context) { gc.JSON(200, gin.H{}) })
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/x", nil)
		req.Header.Set("X-User-ID", "not-a-uuid")
		req.Header.Set("Authorization", "Bearer "+sign(t, "USER", false))
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatalf("want 200 got %d", w.Code)
		}
	})
}

func TestGuards(t *testing.T) {
	c := newClient(t)

	run := func(t *testing.T, guard gin.HandlerFunc, role string, owner bool) int {
		r := gin.New()
		r.GET("/x", WithAuth(c, WithLocalValidation()), guard, func(gc *gin.Context) { gc.JSON(200, gin.H{}) })
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/x", nil)
		req.Header.Set("Authorization", "Bearer "+sign(t, role, owner))
		r.ServeHTTP(w, req)
		return w.Code
	}

	if run(t, OwnerOnly(), "USER", true) != 200 {
		t.Error("owner should pass")
	}
	if run(t, OwnerOnly(), "USER", false) != 403 {
		t.Error("non-owner should 403")
	}
	if run(t, AdminOnly(), "ADMIN", false) != 200 {
		t.Error("admin should pass")
	}
	if run(t, AdminOnly(), "USER", false) != 403 {
		t.Error("non-admin should 403")
	}
	if run(t, RequireRole(auth.RoleAdmin, auth.RoleUser), "USER", false) != 200 {
		t.Error("USER in role set should pass")
	}
}

func TestRequireServiceToken(t *testing.T) {
	c, _ := auth.New(auth.Config{InternalSecret: "internal"})
	r := gin.New()
	r.GET("/svc", RequireServiceToken(c), func(gc *gin.Context) { gc.JSON(200, gin.H{}) })

	tok, err := c.GenerateServiceToken()
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/svc", nil)
	req.Header.Set("X-Service-Token", tok)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("valid svc token want 200 got %d", w.Code)
	}

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/svc", nil)
	r.ServeHTTP(w2, req2)
	if w2.Code != 401 {
		t.Fatalf("missing svc token want 401 got %d", w2.Code)
	}
}
