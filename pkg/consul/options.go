package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

func WithServiceCheck(addr string, port int) Option {
	return newOptFunc(func(consul *Consul) {
		var i, t, d string

		if consul.service.interval != "" {
			i = consul.service.interval
		} else {
			i = "10s"
		}

		if consul.service.timeout != "" {
			t = consul.service.timeout
		} else {
			t = "5s"
		}

		if consul.service.deregisterTimeout != "" {
			d = consul.service.deregisterTimeout
		} else {
			d = "30s"
		}

		consul.service.check = &api.AgentServiceCheck{
			Name:                           fmt.Sprintf("%s-%d", addr, port),
			HTTP:                           fmt.Sprintf("http://%s:%d/health", addr, port),
			Interval:                       i,
			Timeout:                        t,
			DeregisterCriticalServiceAfter: d,
		}
	})
}

func WithTag(tag string) Option {
	return newOptFunc(func(consul *Consul) {
		consul.service.tags = append(consul.service.tags, tag)
	})
}

func WithCheckInterval(interval string) Option {
	return newOptFunc(func(consul *Consul) {
		consul.service.interval = interval
	})
}

func WithCheckTimeout(timeout string) Option {
	return newOptFunc(func(consul *Consul) {
		consul.service.timeout = timeout
	})
}

func WithCheckDeregisterTimeout(timeout string) Option {
	return newOptFunc(func(consul *Consul) {
		consul.service.deregisterTimeout = timeout
	})
}

func WithCheckTLL(timeout string) Option {
	return newOptFunc(func(consul *Consul) {
		consul.service.tll = timeout
	})
}

func WithSelfCheckTimeout(timeout time.Duration) Option {
	return newOptFunc(func(consul *Consul) {
		consul.service.agentSelfTimeout = timeout
	})
}

var _ Option = (*funcOption)(nil)

type Option interface {
	apply(consul *Consul)
}

type funcOption struct {
	f func(consul *Consul)
}

func (fdo *funcOption) apply(consul *Consul) {
	fdo.f(consul)
}

func newOptFunc(f func(consul *Consul)) *funcOption {
	return &funcOption{f: f}
}
