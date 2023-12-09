package container

import (
	"encoding/json"
	"os"

	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
)

// AppConfig
type AppConfig struct {
	Env              string   `json:"env"`
	EnvList          []string `env:"ENV_LIST" envSeparator:","`
	RouterCategory   []string `env:"ROUTER_CATEGORY" envSeparator:","`
	ConsolePort      int      `env:"CONSOLE_PORT"`
	StaticDir        string   `env:"STATIC_DIR"`
	Port             int64    `env:"PORT"`
	Database         string   `env:"DATABASE"`
	LogDir           string   `env:"LOG_DIR"`
	SyncInterval     int64    `env:"SYNC_INTERVAL"`
	Debug            bool     `env:"DEBUG"`
	LoginSalt        string   `env:"LOG_SALT"`
	PermissionConfig string   `env:"PERMISSION_CONFIG"`
}

// Permission
type Permission struct {
	Roles []struct {
		Role        string `json:"role"`
		Name        string `json:"name"`
		Superuser   bool   `json:"superuser"`
		Description string `json:"description"`
	} `json:"roles"`
	Permissions map[string]string `json:"permissions"`
}

var (
	// AppConfig ...
	config *AppConfig

	permission *Permission = &Permission{}
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
	if len(config.PermissionConfig) > 0 {
		if bytes, err := os.ReadFile(cfg.PermissionConfig); err == nil {
			json.Unmarshal(bytes, permission)
		}
	}
	return config, nil
}

// GetAppConfig ...
//
//	@return *AppConfig
func GetAppConfig() *AppConfig {
	return config
}

func GetPermission() *Permission {
	return permission
}
