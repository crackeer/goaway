package container

import (
	"errors"
	"fmt"
	"strings"
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
	data, errorList := getRouter(config.RouterDir, config.DBConnection)
	if len(errorList) > 0 {
		panic(fmt.Sprintf("get routers error:%s", strings.Join(errorList, ";")))
	}
	for path, c := range data {
		routerCache.Set(path, c, cache.DefaultExpiration)
	}
	return nil
}

func getRouter(routerDir, dbConnection string) (map[string]*router.RouterConfig, []string) {
	var (
		retData   map[string]*router.RouterConfig = map[string]*router.RouterConfig{}
		errorList []string
	)
	if len(routerDir) > 0 {
		tmp, err := router.GetRouterFromLocal(routerDir)
		if err != nil {
			errorList = append(errorList, err.Error())
		} else {
			for key, value := range tmp {
				retData[key] = value
			}
		}
	}

	if len(dbConnection) > 0 {
		db, err := OpenDatabase(config.DBConnection)
		if err != nil {
			errorList = append(errorList, fmt.Sprintf("connect %s error:%s", config.DBConnection, err.Error()))
		} else {
			tmp, err := router.GetRouterFromDB(db)
			if err != nil {
				panic(err.Error())
			}
			for key, value := range tmp {
				retData[key] = value
			}
		}
	}
	return retData, errorList
}
