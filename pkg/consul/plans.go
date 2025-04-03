package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/labstack/gommon/log"
)

const planType = "service"

type Plan struct {
	consulURL string
	service   string
	plan      *watch.Plan
	input     chan<- []*api.ServiceEntry
	errCh     chan<- error
}

func NewPlan(consulURL string, serviceName string, input chan<- []*api.ServiceEntry) *Plan {
	var p = &Plan{}

	pl, _ := watch.Parse(map[string]interface{}{
		"type":        planType,
		"service":     serviceName,
		"passingonly": true,
	})

	log.Debug("new consul plan", serviceName)

	p.consulURL = consulURL
	p.service = serviceName
	p.plan = pl
	p.input = input

	pl.Handler = p.handle

	return p
}

func (p *Plan) handle(_ uint64, data interface{}) {
	if !p.plan.IsStopped() {
		entries := data.([]*api.ServiceEntry)
		if entries != nil && len(entries) > 0 {
			p.input <- entries
		}
	}
}

func (p *Plan) Run(errCh chan<- error) {
	go func() {
		if err := p.plan.Run(p.consulURL); err != nil {
			errCh <- err
		}
	}()

	p.errCh = errCh
}

func (p *Plan) Stop() {
	p.plan.Stop()

	if p.plan.IsStopped() {
		close(p.input)
		close(p.errCh)
	}
}
