package config

import "github.com/caarlos0/env/v6"

type Config struct {
	GRPCPort int `env:"GRPC_PORT" envDefault:"9090"`
	Port     int `env:"PORT" envDefault:"8080"`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
