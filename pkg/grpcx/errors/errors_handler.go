package errors

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	apperrors "github.com/Brain-Wave-Ecosystem/go-common/pkg/error"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HTTPError struct {
	Message    string `json:"message"`
	IncidentID string `json:"incident_id,omitempty"`
}

func NewCustomErrorHandler(logger *zap.Logger) func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	return func(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		var appError *apperrors.Error
		if !errors.As(err, &appError) {
			appError = apperrors.InternalWithoutStackTrace(err)
		}

		grpcErr, found := status.FromError(err)
		if !found {
			grpcErr = status.Convert(err)
		}

		grpcCode := grpcErr.Code()

		httpError := HTTPError{
			IncidentID: appError.IncidentID,
			Message:    grpcErr.Message(),
		}

		if grpcCode == codes.Internal {
			logger.Error("internal error",
				zap.String("code", grpcCode.String()),
				zap.String("incident_id", appError.IncidentID),
				zap.String("message", grpcErr.Message()),
				zap.String("method", r.Method),
				zap.String("url", r.RequestURI),
				zap.String("stack_trace", appError.StackTrace),
			)
		} else {
			logger.Debug("client error", zap.String("message", grpcErr.Message()))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(runtime.HTTPStatusFromCode(grpcCode))
		_ = json.NewEncoder(w).Encode(httpError)
	}
}
