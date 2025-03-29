package middlewares

import (
	"fmt"
	"net/http"

	"github.com/Brain-Wave-Ecosystem/go-common/pkg/telemetry/helpers"

	"go.opentelemetry.io/otel"
)

func TraceTelemetryMiddleware(tracerName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tracer := otel.Tracer(tracerName)
			ctx, span := tracer.Start(r.Context(), fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path))
			defer span.End()

			ctx = helpers.InjectTracerToMetadata(ctx)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
