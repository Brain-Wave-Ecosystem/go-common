package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

func WithServiceCheck(addr string, port int) Option {
	return newOptFunc(func(consul *Consul) {
		var i, t, d string

		if consul.interval != "" {
			i = consul.interval
		} else {
			i = "10s"
		}

		if consul.timeout != "" {
			t = consul.timeout
		} else {
			t = "5s"
		}

		if consul.deregisterTimeout != "" {
			d = consul.deregisterTimeout
		} else {
			d = "30s"
		}

		consul.check = &api.AgentServiceCheck{
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
		consul.tag = tag
	})
}

func WithCheckInterval(interval string) Option {
	return newOptFunc(func(consul *Consul) {
		consul.interval = interval
	})
}

func WithCheckTimeout(timeout string) Option {
	return newOptFunc(func(consul *Consul) {
		consul.timeout = timeout
	})
}

func WithCheckDeregisterTimeout(timeout string) Option {
	return newOptFunc(func(consul *Consul) {
		consul.deregisterTimeout = timeout
	})
}

func WithCheckTLL(timeout string) Option {
	return newOptFunc(func(consul *Consul) {
		consul.tll = timeout
	})
}

func WithSelfCheckTimeout(timeout time.Duration) Option {
	return newOptFunc(func(consul *Consul) {
		consul.agentSelfTimeout = timeout
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
