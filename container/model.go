package container

import (
	"errors"
	"fmt"
	"time"

	"github.com/crackeer/go-gateway/model"
	"github.com/crackeer/go-gateway/util/database"
	apiBase "github.com/crackeer/simple_http"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

var (
	routerCache *cache.Cache
	appCache    *cache.Cache
	modelDB     *gorm.DB
)

func init() {
	routerCache = cache.New(20*time.Minute, 30*time.Minute)
	appCache = cache.New(20*time.Minute, 30*time.Minute)
}

// InitModel
func InitModel() {
	if db, err := database.Open(config.DBConnection); err != nil {
		panic(fmt.Sprintf("open database `%s`: %v", config.DBConnection, err))
	} else {
		modelDB = db
	}

}

func saveModel() {
	routers, _ := model.GetRouterFromDB(modelDB)
	for path, value := range routers {
		routerCache.Set(path, value, cache.DefaultExpiration)
	}

	apis, _ := model.GetServiceAPIFromDB(modelDB)
	for name, value := range apis {
		apiBase.RegisterServiceAPI(name, value)
	}

	signs, _ := model.GetSignCodeFromDB(modelDB)
	for key, value := range signs {
		err := apiBase.RegisterLuaSign(key, value)
		if err != nil {
			//panic(fmt.Errorf("register sign error: %v[%s]", err.Error(), key))
		}
	}
	apps, _ := model.GetAppList(modelDB)
	for path, value := range apps {
		routerCache.Set(path, value, cache.DefaultExpiration)
	}

}

// GetRouter
//
//	@param path
//	@return *model.RouterConfig
//	@return error
func GetRouter(path string) (*model.RouterConfig, error) {
	if value, exists := routerCache.Get(path); exists {
		return value.(*model.RouterConfig), nil
	}
	return nil, errors.New("not found: " + path)
}

// GetApp
//
//	@param appID
//	@return *model.AppConfig
//	@return error
func GetApp(appID string) (*model.AppConfig, error) {
	if value, exists := routerCache.Get(appID); exists {
		return value.(*model.AppConfig), nil
	}
	return nil, errors.New("not found: " + appID)
}
