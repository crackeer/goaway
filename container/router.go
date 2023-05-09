package container

import (
	"net/http"

	api "github.com/crackeer/simple_http"
)

// InitRouterFactory
//
//	@param cfg
func InitRouterFactory(cfg *AppConfig) error {
	return nil
}

// InitAPI
//
//	@param cfg
//	@return error
func InitAPI(cfg *AppConfig) error {
	api.RegisterServiceAPI("abc/test", &api.ServiceAPI{
		Host:           "https://www.boredapi.com",
		DisableExtract: true,
		SignType:       "test",
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
	}
	return nil
}
