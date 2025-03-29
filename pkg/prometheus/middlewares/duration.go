package middlewares

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusDurationMiddleware(c *prometheus.HistogramVec) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)
			duration := time.Since(start)
			c.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
		})
	}
}
