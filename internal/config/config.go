package config

import "github.com/caarlos0/env/v6"

type Config struct {
	GRPCPort int    `env:"GRPC_PORT" envDefault:"9090"`
	Port     int    `env:"PORT" envDefault:"8080"`
	Host     string `env:"HOST" envDefault:"localhost"`
	User     string `env:"USERDB" envDefault:"postgres"`
	Password string `env:"PASSWORD" envDefault:"postgres"`
	Dbname   string `env:"DBNAME" envDefault:"backend"`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
