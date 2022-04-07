package metrics

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MetricRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "HTTP Request Count.",
		},
		[]string{"status", "path"},
	)
	MetricRequestDurationSecondsBucket = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds_bucket",
			Help:    "Histogram of latencies for HTTP requests.",
			Buckets: []float64{0.1, 0.2, 0.4, 1, 3, 8, 20, 60, 120},
		},
		[]string{"status", "path"},
	)
)

func init() {
	prometheus.MustRegister(MetricRequestsTotal)
	prometheus.MustRegister(MetricRequestDurationSecondsBucket)
}

type Config struct {
	IgnorePaths []string `conf:"ignore_paths"`
}

func MetricsMiddleware(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Check if the prefix should be ignored
			for _, prefix := range config.IgnorePaths {
				if strings.HasPrefix(r.RequestURI, prefix) {
					next.ServeHTTP(w, r)
					return
				}
			}

			start := time.Now()

			// See if we're already using a wrapped response writer and if not make one.
			ww, ok := w.(middleware.WrapResponseWriter)
			if !ok {
				ww = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			}

			next.ServeHTTP(ww, r)

			path := chi.RouteContext(r.Context()).RoutePattern()

			MetricRequestsTotal.With(
				prometheus.Labels{
					"path":   path,
					"status": strconv.Itoa(ww.Status()),
				}).Inc()
			MetricRequestDurationSecondsBucket.With(
				prometheus.Labels{
					"path":   path,
					"status": strconv.Itoa(ww.Status()),
				}).Observe(time.Since(start).Seconds())
		})
	}
}
