package container

import (
	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
)

type AppConfig struct {
	Env          string `env:"ENV"`
	Port         int64  `env:"PORT"`
	DBConnection string `env:"DB_CONNECTION"`
	RouterDir    string `env:"ROUTER_DIR"`
	APIDir       string `env:"API_DIR"`
	SignDir      string `env:"SIGN_DIR"`
	LogDir       string `env:"LOG_DIR"`

	SyncInterval int64 `env:"SYNC_INTERVAL"`
	Debug        bool  `env:"DEBUG"`
}

var (
	// AppConfig ...
	config *AppConfig
)

// InitAppConfig  ...
//
//	@return *AppConfig
//	@return error
func InitAppConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	config = cfg
	return config, nil
}

// GetAppConfig ...
//
//	@return *AppConfig
func GetAppConfig() *AppConfig {
	return config
}
