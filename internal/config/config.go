package config

import "github.com/caarlos0/env/v6"

// Config ...
type Config struct {
	ServiceAddress      string `env:"SERVICE_PORT" envDefault:"0.0.0.0:8083"`
	DatabaseURI         string `env:"DATABASE_URI,file" envDefault:"./env/database_uri"`
	DatabaseName        string `env:"DATABASE_NAME,file" envDefault:"./env/database_name"`
	AccessTokenSecret   string `env:"ACCESS_TOKEN_SECRET" envDefault:"m8KP74IcTMiOYEhFP2Da"`
	AccessTokenLifetime string `env:"ACCESS_TOKEN_LIFETIME" envDefault:"8h"`
}

// NewConfig ...
func NewConfig() *Config {
	c := new(Config)
	if err := env.Parse(c); err != nil {
		panic(err)
	}
	return c
}
