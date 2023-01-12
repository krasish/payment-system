package config

import "github.com/vrischmann/envconfig"

type Config struct {
	HttpConfig
	DatabaseConfig
}

func NewConfigFromEnv() (Config, error) {
	cfg := Config{}
	err := envconfig.Init(&cfg)
	return cfg, err
}
