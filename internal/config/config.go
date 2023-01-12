package config

import "github.com/vrischmann/envconfig"

type Config struct {
	HttpConfig
	DatabaseConfig
	ViewTemplatesPath string `envconfig:"APP_VIEW_TEMPLATES_PATH"`
}

func NewConfigFromEnv() (Config, error) {
	cfg := Config{}
	err := envconfig.Init(&cfg)
	return cfg, err
}
