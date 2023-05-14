package container

import (
	"errors"
	"time"

	"github.com/crackeer/go-gateway/base/router"
	"github.com/patrickmn/go-cache"
)

var (
	routerCache *cache.Cache
)

func init() {
	routerCache = cache.New(20*time.Minute, 30*time.Minute)
}

func GetRouter(path string) (*router.RouterConfig, error) {
	if value, exists := routerCache.Get(path); exists {
		return value.(*router.RouterConfig), nil
	}
	return nil, errors.New("not found: " + path)
}

func InitRouter() error {
	data := map[string]*router.RouterConfig{}
	if len(config.RouterDir) > 0 {
		tmp, err := router.GetRouterFromLocal(config.RouterDir)
		if err != nil {
			panic(err.Error())
		}
		for key, value := range tmp {
			data[key] = value
		}
	}

	if len(config.SqliteFile) > 0 {
		tmp, err := router.GetRouterFromSQLite(config.SqliteFile)
		if err != nil {
			panic(err.Error())
		}
		for key, value := range tmp {
			data[key] = value
		}
	}
	for path, c := range data {
		routerCache.Set(path, c, cache.DefaultExpiration)
	}
	return nil
}
