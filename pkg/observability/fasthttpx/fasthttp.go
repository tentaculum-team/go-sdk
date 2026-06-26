// Package fasthttpx adapts the observability RED metrics to fasthttp, for the
// protect service. Kept separate so gin-only consumers don't pull fasthttp.
package fasthttpx

import (
	"strconv"
	"time"

	"github.com/tentaculum-team/go-sdk/pkg/observability"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// MetricsHandler serves the Prometheus registry over fasthttp. Mount it at
// /metrics. The route label is fixed since fasthttp's router doesn't expose the
// matched template here; per-route timing comes from Middleware below.
func MetricsHandler() fasthttp.RequestHandler {
	return fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
}

// Middleware wraps a handler to record RED metrics. Pass the matched route
// template (e.g. "/apps/{id}") as route to keep label cardinality bounded; pass
// the raw path only if no template is available.
func Middleware(route string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		next(ctx)
		observability.Observe(
			string(ctx.Method()),
			route,
			strconv.Itoa(ctx.Response.StatusCode()),
			time.Since(start).Seconds(),
		)
	}
}
