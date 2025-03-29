package consul

import (
	"fmt"
	"github.com/Brain-Wave-Ecosystem/go-common/pkg/clients"
	"time"

	"go.uber.org/zap"

	apperrors "github.com/Brain-Wave-Ecosystem/go-common/pkg/error"

	"github.com/hashicorp/consul/api"
)

type Consul struct {
	client    *api.Client
	consulURL string
	logger    *zap.Logger
	service   *ServiceConfig
	plans     []*Plan
}

type ServiceConfig struct {
	id                string
	name              string
	addr              string
	tags              []string
	grpcPort          int
	check             *api.AgentServiceCheck
	interval          string
	timeout           string
	tll               string
	deregisterTimeout string
	agentSelfTimeout  time.Duration
}

func NewConsul(consulURL, name, address string, grpcPort int, logger *zap.Logger, options ...Option) (*Consul, error) {
	client, err := clients.NewConsulClient(consulURL)
	if err != nil {
		return nil, err
	}

	var c Consul

	c.client = client
	c.consulURL = consulURL
	c.logger = logger

	c.service = &ServiceConfig{
		name:              name,
		addr:              address,
		grpcPort:          grpcPort,
		tags:              []string{"v1", "http", "grpc"},
		interval:          "10s",
		timeout:           "5s",
		tll:               "15s",
		deregisterTimeout: "30s",
		agentSelfTimeout:  20 * time.Second,
	}

	for _, option := range options {
		option.apply(&c)
	}

	c.service.id = fmt.Sprintf("%s-%d", c.service.name, c.service.grpcPort)

	return &c, nil
}

func (c *Consul) Consul() *api.Client {
	return c.client
}

func (c *Consul) RegisterService() error {
	return c.register()
}

func (c *Consul) WatchService(serviceAddr string) <-chan []*api.ServiceEntry {
	queue := make(chan []*api.ServiceEntry, 3)

	plan := NewPlan(c.consulURL, serviceAddr, queue)

	c.plans = append(c.plans, plan)

	return queue
}

func (c *Consul) Stop() error {
	for _, plan := range c.plans {
		plan.Stop()
	}

	return c.client.Agent().ServiceDeregister(c.service.id)
}

func (c *Consul) register() error {
	registration := &api.AgentServiceRegistration{
		ID:      c.service.id,
		Name:    c.service.name,
		Address: c.service.addr,
		Port:    c.service.grpcPort,
		Tags:    c.service.tags,
	}

	if c.service.check != nil {
		registration.Check = c.service.check
	} else {
		registration.Check = &api.AgentServiceCheck{
			Name:                           fmt.Sprintf("%s-%d", c.service.addr, c.service.grpcPort),
			GRPC:                           fmt.Sprintf("%s:%d", c.service.addr, c.service.grpcPort),
			Interval:                       c.service.interval,
			Timeout:                        c.service.timeout,
			DeregisterCriticalServiceAfter: c.service.deregisterTimeout,
		}
	}

	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		return apperrors.Internal(err)
	}

	c.logger.Info("Service registered in Consul", zap.String("name", c.service.name), zap.String("address", c.service.addr), zap.Strings("tags", c.service.tags))

	return nil
}
