package middlewares

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusRequestSizeSummaryMiddleware(c *prometheus.HistogramVec) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{ResponseWriter: w, statusCode: 200}
			next.ServeHTTP(rw, r)

			length := float64(r.ContentLength)
			if length != -1 {
				status := strconv.Itoa(rw.statusCode)
				c.WithLabelValues(r.Method, r.URL.Path, status).Observe(length)
			}
		})
	}
}

func PrometheusResponseSizeSummaryMiddleware(c *prometheus.HistogramVec) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{ResponseWriter: w, statusCode: 200}
			next.ServeHTTP(rw, r)

			if rw.bytesAmount > 0 {
				status := strconv.Itoa(rw.statusCode)
				c.WithLabelValues(r.Method, r.URL.Path, status).Observe(float64(rw.bytesAmount))
			}
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	bytesAmount int
	statusCode  int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	amount, err := rw.ResponseWriter.Write(b)
	rw.bytesAmount = amount
	return amount, err
}
