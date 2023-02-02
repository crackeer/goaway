package container

import (
	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	Env  string `env:"ENV"`
	Port int64  `env:"PORT"`
}

var (
	// AppConfig ...
	Config *config
)

// InitConfig ...
func InitConfig() error {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		return err
	}
	Config = cfg
	return nil
}
