package container

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/crackeer/goaway/model"
	apiBase "github.com/crackeer/simple_http"
	"github.com/glebarez/sqlite"
	"github.com/patrickmn/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	routerCache *cache.Cache
	modelDB     *gorm.DB
)

func init() {
	routerCache = cache.New(20*time.Minute, 30*time.Minute)
}

func GetModelDB() *gorm.DB {
	return modelDB
}

// InitModel
func InitModel() {
	if db, err := open(config.Database); err != nil {
		panic(fmt.Sprintf("open database `%s`: %v", config.Database, err))
	} else {
		modelDB = db
	}
}
func open(connection string) (*gorm.DB, error) {
	if strings.HasPrefix(connection, "mysql://") {
		return gorm.Open(mysql.Open(connection[8:]), &gorm.Config{})
	}

	if strings.HasPrefix(connection, "sqlite://") {
		return gorm.Open(sqlite.Open(connection[9:]), &gorm.Config{})
	}

	return nil, errors.New("not support")
}

func saveModel() {
	routers, _ := model.GetRouterFromDB(modelDB)
	for path, value := range routers {
		routerCache.Set(path, value, cache.DefaultExpiration)
	}

	apis, _ := model.GetServiceAPIFromDB(modelDB, config.Env)
	for name, value := range apis {
		apiBase.RegisterServiceAPI(name, value)
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
