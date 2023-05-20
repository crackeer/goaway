package container

import (
	"fmt"

	apiBase "github.com/crackeer/simple_http"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
)

var cronTab *cron.Cron

// StartSchedule
//
//	@param ctx
//	@param cfg
func StartSchedule() {

	cronTab = cron.New(cron.WithSeconds())
	cronTab.AddFunc("1 * * * * *", func() {
		fmt.Println("crontab")
		signList, errorsList := getSign(config.SignDir, config.DBConnection)
		if len(errorsList) > 0 {
			Log(map[string]interface{}{
				"errors": errorsList,
			}, "SyncSignError")
		}
		for key, value := range signList {
			err := apiBase.RegisterLuaSign(key, value)
			if err != nil {
				Log(map[string]interface{}{
					"key":   key,
					"value": value,
				}, "RegisterLuaSignError")
			}
		}
		apiList, errorsList := getServiceAPIMap(config.APIDir, config.DBConnection)
		if len(errorsList) > 0 {
			Log(map[string]interface{}{
				"errors": errorsList,
			}, "SyncAPIError")
		}
		for name, c := range apiList {
			apiBase.RegisterServiceAPI(name, c)
		}

		routers, errorsList := getRouter(config.RouterDir, config.DBConnection)
		if len(errorsList) > 0 {
			Log(map[string]interface{}{
				"errors": errorsList,
			}, "SyncRouterError")
		}
		for path, c := range routers {
			routerCache.Set(path, c, cache.DefaultExpiration)
		}
		Log(map[string]interface{}{}, "CrontabFinished")
	})
	cronTab.Start()
}
