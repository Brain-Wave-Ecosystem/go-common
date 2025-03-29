package interceptors

import (
	"context"
	"fmt"

	"github.com/Brain-Wave-Ecosystem/go-common/pkg/telemetry/helpers"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

func TraceTelemetryServerInterceptor(tracerName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		tracer := otel.Tracer(fmt.Sprintf("server: %s", tracerName))

		ctx = helpers.ExtractTracerFromMetadata(ctx)
		ctx, span := tracer.Start(ctx, info.FullMethod)
		defer span.End()

		return handler(ctx, req)
	}
}

func TraceTelemetryClientInterceptor(tracerName string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		tracer := otel.Tracer(fmt.Sprintf("client: %s", tracerName))

		ctx, span := tracer.Start(ctx, method)
		defer span.End()

		ctx = helpers.InjectTracerToMetadata(ctx)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
