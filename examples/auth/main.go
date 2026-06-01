// Example: protect a Gin API with the auth SDK.
//
//	AUTH_ENV=dev AUTH_URL_DEV=http://localhost:8080 go run ./examples/auth
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Tentaculum-dev/go-sdk/auth"
	"github.com/Tentaculum-dev/go-sdk/auth/cache"
	ginmw "github.com/Tentaculum-dev/go-sdk/auth/middleware/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	// Config from env: picks AUTH_URL_DEV / AUTH_URL_PROD by AUTH_ENV,
	// falls back to AUTH_URL. JWT_SECRET (offline) + INTERNAL_JWT_SECRET optional.
	cfg := auth.ConfigFromEnv()
	if cfg.BaseURL == "" {
		log.Fatal("set AUTH_URL (or AUTH_URL_DEV/PROD)")
	}
	// Opt-in remote-validation cache (60s cap by default).
	cfg.Cache = cache.NewLRU(1024)

	client, err := auth.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	// All /api/v1 routes require a valid bearer token (remote validation).
	api := r.Group("/api/v1", ginmw.WithAuth(client))
	api.GET("/things", func(c *gin.Context) {
		id, _ := ginmw.IdentityFrom(c)
		c.JSON(http.StatusOK, gin.H{
			"user_id": id.UserID, "org_id": id.OrgID, "role": id.Role,
		})
	})

	// Owner-only subgroup.
	owner := api.Group("/admin", ginmw.OwnerOnly())
	owner.POST("/things", func(c *gin.Context) { c.JSON(http.StatusCreated, gin.H{"ok": true}) })

	// Service-to-service endpoint (X-Service-Token), if INTERNAL_JWT_SECRET set.
	r.GET("/internal/ping", ginmw.RequireServiceToken(client), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pong": true})
	})

	addr := ":" + envOr("PORT", "9090")
	log.Printf("listening on %s (prod=%v)", addr, auth.IsProd())
	log.Fatal(r.Run(addr))
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
