package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type ServiceConfig interface {
	SetDefaults()
}

type DefaultConfig struct {
	Name                    string `env:"NAME"`
	Address                 string `env:"ADDRESS"`
	Tag                     string `env:"TAG" envDefault:"v1"`
	Local                   bool   `env:"LOCAL" envDefault:"false"`
	WebPort                 int    `env:"PORT" envDefault:"8010"`
	GRPCPort                int    `env:"GRPC_PORT" envDefault:"8011"`
	GracefulShutdownTimeout int    `env:"GRACEFUL_TIMEOUT" envDefault:"10"`
	ConsulAddr              string `env:"CONSUL_ADDRESS" envDefault:"127.0.0.1:8500"`
	TracerAddr              string `env:"TRACER_ADDRESS" envDefault:"http://jaeger:14268/api/traces"`
}

func NewConfig[T any]() (*T, error) {
	_ = godotenv.Load()

	config := new(T)
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return config, nil
}
