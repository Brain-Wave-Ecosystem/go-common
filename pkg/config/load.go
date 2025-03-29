package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Load[C any]() (*C, error) {
	_ = godotenv.Load()

	cfg := new(C)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}
