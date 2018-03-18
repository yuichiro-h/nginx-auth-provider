package config

import "github.com/kelseyhightower/envconfig"

const envPrefix = "nginx_auth_provider"

var c Config

type Config struct {
	Port  int  `envconfig:"port"`
	Debug bool `envconfig:"debug"`

	CookieSecret string `envconfig:"cookie_secret"`
	CookieMaxAge int    `envconfig:"cookie_max_age"`

	GoogleDomain       string `envconfig:"google_domain"`
	GoogleClientID     string `envconfig:"google_client_id"`
	GoogleClientSecret string `envconfig:"google_client_secret"`
	GoogleCallbackURL  string `envconfig:"google_callback_url"`
}

func Load() {
	envconfig.MustProcess(envPrefix, &c)
}

func Get() *Config {
	return &c
}
