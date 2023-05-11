package container

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	routerCache *cache.Cache
)

// RouterConfig
type RouterConfig struct {
	Mode     string      `json:"mode"`
	ProxyAPI string      `json:"proxy_api"`
	Response interface{} `json:"response"`
}

func init() {
	routerCache = cache.New(20*time.Minute, 30*time.Minute)
}

// GetRouterConfig
//
//	@param path
//	@return *RouterConfig
//	@return error
func GetRouterConfig(path string) (*RouterConfig, error) {
	if value, exists := routerCache.Get(path); exists {
		return value.(*RouterConfig), nil
	}
	c, err := readRouterConfig(path)
	if err != nil {
		return nil, err
	}

	routerCache.Set(path, c, cache.DefaultExpiration)
	return c, nil
}

func readRouterConfig(path string) (*RouterConfig, error) {
	fullPath := filepath.Join(GetAppConfig().RouterDir, path)
	bytes, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("read file `%s` error:%s", fullPath, err.Error())
	}
	retData := &RouterConfig{}
	if err := json.Unmarshal(bytes, retData); err != nil {
		return nil, fmt.Errorf("json unmarshal `%s` content error:%s", fullPath, err.Error())
	}
	return retData, nil
}
