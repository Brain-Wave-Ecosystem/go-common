package resolver

import (
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
)

const customScheme = "dynamic"

var _ resolver.Builder = (*Builder)(nil)

func NewBuilder(output <-chan []*api.ServiceEntry, logger *zap.Logger) *Builder {
	return &Builder{
		output: output,
		logger: logger,
	}
}

type Builder struct {
	output <-chan []*api.ServiceEntry
	logger *zap.Logger
}

func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := NewResolver(target, cc, opts, b.output, b.logger)

	go r.watch()

	return r, nil
}

func (b *Builder) Scheme() string {
	return customScheme
}
