package container

import (
	"github.com/crackeer/go-gateway/base/api"
	apiBase "github.com/crackeer/simple_http"
)

// InitAPI
//
//	@param cfg
//	@return error
func InitAPI() error {
	apiList := map[string]*apiBase.ServiceAPI{}
	if len(config.APIDir) > 0 {
		list, err := api.GetServiceAPIFromDir(config.APIDir)
		if err != nil {
			panic(err.Error())
		}
		for key, value := range list {
			apiList[key] = value
		}
	}

	if len(config.SqliteFile) > 0 {
		list, err := api.GetServiceAPIFromSQLite(config.SqliteFile)
		if err != nil {
			panic(err.Error())
		}
		for key, value := range list {
			apiList[key] = value
		}
	}

	for name, c := range apiList {
		apiBase.RegisterServiceAPI(name, c)
	}

	/*
		if len(config.SignDir) > 0 {
			if err := registerSign(config.SignDir); err != nil {
				panic(err.Error())
			}
		}
	*/

	/*
			api.RegisterServiceAPI("abc/test", &api.ServiceAPI{
				Host:           "https://www.boredapi.com",
				DisableExtract: true,
				SignName:       "test",
				SignConfig: map[string]interface{}{
					"ak": "22",
				},
				Path:    "api/activity",
				Method:  http.MethodGet,
				Timeout: 3000,
			})
		err := api.RegisterLuaSignByFile("test", "./config/sign/test.lua")
		if err != nil {
			panic(err.Error())
		}*/
	return nil
}
