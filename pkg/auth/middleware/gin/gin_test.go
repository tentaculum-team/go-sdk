package ginmw

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/Tentaculum-dev/go-sdk/pkg/auth"
	"github.com/gin-gonic/gin"
)

var testSecret = paseto.NewV4AsymmetricSecretKey()

func init() { gin.SetMode(gin.TestMode) }

func sign(t *testing.T, role string) string {
	t.Helper()
	tok := paseto.NewToken()
	now := time.Now()
	tok.SetIssuedAt(now)
	tok.SetExpiration(now.Add(time.Minute))
	tok.SetString("typ", "access")
	tok.SetString("user_uuid", "u1")
	tok.SetString("sys_role", role)
	return tok.V4Sign(testSecret, nil)
}

func newClient(t *testing.T) *auth.Client {
	c, err := auth.New(auth.Config{AccessPublicKey: testSecret.Public().ExportHex()})
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
		gc.JSON(200, gin.H{"role": id.SysRole})
	})

	t.Run("valid bearer", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/x", nil)
		req.Header.Set("Authorization", "Bearer "+sign(t, "ADMIN"))
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

func TestRequireRole(t *testing.T) {
	c := newClient(t)

	run := func(t *testing.T, role string, allowed ...string) int {
		r := gin.New()
		r.GET("/x", WithAuth(c, WithLocalValidation()), RequireRole(allowed...), func(gc *gin.Context) { gc.JSON(200, gin.H{}) })
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/x", nil)
		req.Header.Set("Authorization", "Bearer "+sign(t, role))
		r.ServeHTTP(w, req)
		return w.Code
	}

	if run(t, "ADMIN", auth.RoleAdmin) != 200 {
		t.Error("admin should pass")
	}
	if run(t, "USER", auth.RoleAdmin) != 403 {
		t.Error("non-admin should 403")
	}
	if run(t, "USER", auth.RoleAdmin, auth.RoleUser) != 200 {
		t.Error("USER in role set should pass")
	}
}
