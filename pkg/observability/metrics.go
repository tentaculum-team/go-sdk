// Package observability gives every Tentaculum service the same Prometheus
// metrics with one import: RED metrics (Rate, Errors, Duration) for HTTP plus
// the default Go runtime/process collectors, exposed on /metrics.
//
// Wiring (gin):
//
//	engine.Use(observability.Middleware())
//	observability.Register(engine) // adds GET /metrics
//
// Wiring (fasthttp): see the fasthttpx subpackage.
//
// All series go on the default Prometheus registry, so promhttp.Handler()
// already serves them alongside go_* and process_* metrics.
package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests by method, route and status code.",
	}, []string{"method", "route", "status"})

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency in seconds by method, route and status code.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route", "status"})
)

// Observe records one finished HTTP request. Framework adapters (gin, fasthttp)
// call this; route should be the matched template, not the raw path, to keep
// label cardinality bounded.
func Observe(method, route, status string, seconds float64) {
	if route == "" {
		route = "unmatched"
	}
	httpRequests.WithLabelValues(method, route, status).Inc()
	httpDuration.WithLabelValues(method, route, status).Observe(seconds)
}
