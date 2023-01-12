package config

import "time"

type HttpConfig struct {
	JwtKey          string        `envconfig:"default=secretKey,APP_HTTP_JWT_KEY"`
	TransactionPath string        `envconfig:"default=/transaction,APP_HTTP_TRANSACTION_PATH"`
	MerchantPath    string        `envconfig:"default=/merchant,APP_HTTP_MERCHANT_PATH"`
	UserPath        string        `envconfig:"default=/user,APP_HTTP_USER_PATH"`
	ViewsPath       string        `envconfig:"default=/views,APP_HTTP_VIEWS_PATH"`
	Port            string        `envconfig:"default=8080,APP_HTTP_PORT"`
	ServerTimeout   time.Duration `envconfig:"default=110s,APP_HTTP_SERVER_TIMEOUT"`
}
