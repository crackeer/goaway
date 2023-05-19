package container

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitContainer
//
//	@param ctx
//	@param wg
//	@return error
func InitContainer(ctx context.Context, wg *sync.WaitGroup) error {
	appConfig, err := InitAppConfig()
	if err != nil {
		return fmt.Errorf("init app config failed: %s", err.Error())
	}
	err = InitLogger(appConfig)
	if err != nil {
		return fmt.Errorf("init logger failed: %s", err.Error())
	}
	Log(map[string]interface{}{
		"app_config": appConfig,
	}, "AppConfig")
	InitAPI()
	InitSign()
	InitRouter()
	go StartSchedule(ctx, wg, appConfig)

	return nil
}

// OpenDatabase
//
//	@param connection
//	@return *gorm.DB
//	@return error
func OpenDatabase(connection string) (*gorm.DB, error) {
	if strings.HasPrefix(connection, "mysql://") {
		return gorm.Open(mysql.Open(connection[8:]), &gorm.Config{})
	}

	if strings.HasPrefix(connection, "sqlite://") {
		return gorm.Open(sqlite.Open(connection[9:]), &gorm.Config{})
	}

	return nil, errors.New("not support")
}
