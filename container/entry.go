package container

import (
	"context"
	"fmt"
	"sync"
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
	InitAPI(appConfig)
	go StartSchedule(ctx, wg, appConfig)

	return nil
}
