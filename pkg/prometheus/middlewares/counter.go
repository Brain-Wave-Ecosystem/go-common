package middlewares

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusCounterMiddleware(c *prometheus.CounterVec) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{ResponseWriter: w, statusCode: 200}
			next.ServeHTTP(rw, r)

			status := strconv.Itoa(rw.statusCode)
			c.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		})
	}
}
