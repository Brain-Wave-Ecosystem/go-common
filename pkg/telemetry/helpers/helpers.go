package helpers

import (
	"context"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/metadata"
)

func InjectTracerToMetadata(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, metadataCarrier{md})

	return metadata.NewOutgoingContext(ctx, md)
}

func ExtractTracerFromMetadata(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	propagator := otel.GetTextMapPropagator()
	return propagator.Extract(ctx, metadataCarrier{md})
}

type metadataCarrier struct {
	metadata.MD
}

func (m metadataCarrier) Get(key string) string {
	values := m.MD.Get(key)
	if len(values) > 0 {
		return values[0]
	}

	return ""
}

func (m metadataCarrier) Set(key, value string) {
	m.MD.Set(key, value)
}

func (m metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(m.MD))
	for k := range m.MD {
		keys = append(keys, k)
	}

	return keys
}
