package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"go.uber.org/zap"
)

const planType = "service"

type Plan struct {
	consulURL string
	service   string
	*watch.Plan
	input chan<- []*api.ServiceEntry
	errCh chan<- error

	logger *zap.Logger
}

func NewPlan(consulURL string, serviceName string, input chan<- []*api.ServiceEntry, logger *zap.Logger) *Plan {
	p := &Plan{}

	pl, _ := watch.Parse(map[string]interface{}{
		"type":        planType,
		"service":     serviceName,
		"passingonly": true,
	})

	p.consulURL = consulURL
	p.service = serviceName
	p.Plan = pl
	p.input = input
	p.logger = logger

	pl.Handler = p.handle

	return p
}

func (p *Plan) handle(_ uint64, data interface{}) {
	if !p.IsStopped() {
		p.logger.Debug("Plan working [1]", zap.String("address", p.service))
		entries := data.([]*api.ServiceEntry)
		if entries != nil && len(entries) > 0 {
			p.logger.Debug("Plan working [2]", zap.String("address", p.service))
			p.input <- entries
		}
	}
}

func (p *Plan) Run(errCh chan<- error) {
	go func() {
		if err := p.Plan.Run(p.consulURL); err != nil {
			errCh <- err
		}
	}()

	p.errCh = errCh
}

func (p *Plan) Stop() {
	p.Plan.Stop()

	if p.Plan.IsStopped() {
		close(p.input)
		close(p.errCh)
	}
}
