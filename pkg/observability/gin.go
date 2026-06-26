package observability

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Middleware records RED metrics for every request, labelled by the matched
// route template (c.FullPath), method and status code.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		Observe(
			c.Request.Method,
			c.FullPath(),
			strconv.Itoa(c.Writer.Status()),
			time.Since(start).Seconds(),
		)
	}
}

// Register mounts GET /metrics on the engine.
func Register(engine *gin.Engine) {
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
