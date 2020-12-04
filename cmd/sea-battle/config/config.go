package config

import (
	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Address string `env:"SERVER_ADDRESS"`
}

type Config struct {
	ServerConfig
}

func New() *Config {
	return &Config{
		ServerConfig: ServerConfig{
			Address: "0.0.0.0:8080",
		},
	}
}

func (cfg *Config) Parse() (err error) {
	if err = env.Parse(&cfg.ServerConfig); err != nil {
		return err
	}

	return nil
}
