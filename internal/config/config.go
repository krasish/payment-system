package config

import (
	"time"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	HttpConfig
	DatabaseConfig
	ViewTemplatesPath   string        `envconfig:"APP_VIEW_TEMPLATES_PATH"`
	DeletionJobInterval time.Duration `envconfig:"default=3s,APP_DELETION_JOB_INTERVAL"`
	AdminsImportPath    string        `envconfig:"APP_ADMINS_IMPORT_PATH,optional"`
	MerchantsImportPath string        `envconfig:"APP_MERCHANTS_IMPORT_PATH,optional"`
}

func NewConfigFromEnv() (Config, error) {
	cfg := Config{}
	err := envconfig.Init(&cfg)
	return cfg, err
}
