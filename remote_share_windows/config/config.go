package config

var C = NewAppConfig()

type AppConfig struct {
	Addr string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}
