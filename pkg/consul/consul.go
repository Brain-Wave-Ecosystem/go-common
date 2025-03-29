package consul

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	apperrors "github.com/Brain-Wave-Ecosystem/go-common/pkg/error"

	"github.com/hashicorp/consul/api"
)

type Consul struct {
	*api.Client
	logger            *zap.Logger
	id                string
	name              string
	addr              string
	tag               string
	webPort           int
	grpcPort          int
	check             *api.AgentServiceCheck
	interval          string
	timeout           string
	tll               string
	deregisterTimeout string
	timer             *time.Timer
	agentSelfTimeout  time.Duration
}

func NewConsul(client *api.Client, name, address string, webPort, grpcPort int, logger *zap.Logger, options ...Option) *Consul {
	c := &Consul{Client: client}

	c.name = name
	c.addr = address
	c.webPort = webPort
	c.grpcPort = grpcPort
	c.logger = logger
	c.tag = "v1"
	c.interval = "10s"
	c.timeout = "5s"
	c.tll = "15s"
	c.deregisterTimeout = "30s"
	c.agentSelfTimeout = time.Second * 20

	for _, option := range options {
		option.apply(c)
	}

	c.id = fmt.Sprintf("%s-%s-%d", c.name, c.tag, c.grpcPort)
	c.timer = time.NewTimer(c.agentSelfTimeout)

	return c
}

func (c *Consul) Consul() *api.Client {
	return c.Client
}

func (c *Consul) RegisterService() error {
	return c.register()
}

func (c *Consul) Stop() error {
	c.timer.Stop()
	return c.Client.Agent().ServiceDeregister(c.id)
}

func (c *Consul) register() error {
	registration := &api.AgentServiceRegistration{
		ID:      c.id,
		Name:    c.name,
		Address: c.addr,
		Port:    c.grpcPort,
		Tags:    []string{c.tag},
	}

	if c.check != nil {
		registration.Check = c.check
	} else {
		registration.Check = &api.AgentServiceCheck{
			Name:                           fmt.Sprintf("%s-%d", c.addr, c.webPort),
			HTTP:                           fmt.Sprintf("http://%s:%d/health", c.addr, c.webPort),
			Interval:                       c.interval,
			Timeout:                        c.timeout,
			DeregisterCriticalServiceAfter: c.deregisterTimeout,
		}
	}

	if err := c.Agent().ServiceRegister(registration); err != nil {
		return apperrors.Internal(err)
	}

	c.logger.Info("Service registered in Consul", zap.String("name", c.name), zap.String("address", c.addr), zap.String("tags", c.tag))

	return nil
}
