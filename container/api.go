package container

import (
	"fmt"
	"strings"

	"github.com/crackeer/go-gateway/base/api"
	apiBase "github.com/crackeer/simple_http"
)

// InitAPI
//
//	@param cfg
//	@return error
func InitAPI() error {

	apiList, errorList := getServiceAPIMap(config.APIDir, config.DBConnection)
	if len(errorList) > 0 {
		panic(fmt.Sprintf("get api list error:%s", strings.Join(errorList, ";")))
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

// getServiceAPIMap
//
//	@param cfg
//	@return map[string]*apiBase.ServiceAPI
//	@return []error
func getServiceAPIMap(apiDir string, dbConnection string) (map[string]*apiBase.ServiceAPI, []string) {
	var (
		apiList   map[string]*apiBase.ServiceAPI = map[string]*apiBase.ServiceAPI{}
		errorList []string
	)

	if len(apiDir) > 0 {
		list, err := api.GetServiceAPIFromDir(apiDir)
		if err != nil {
			errorList = append(errorList, err.Error())
		} else {
			for key, value := range list {
				apiList[key] = value
			}
		}
	}

	if len(config.DBConnection) > 0 {
		db, err := OpenDatabase(config.DBConnection)
		if err != nil {
			errorList = append(errorList, fmt.Sprintf("connect %s error:%s", config.DBConnection, err.Error()))
		} else {
			list, err := api.GetServiceAPIFromDB(db)
			if err != nil {
				errorList = append(errorList, err.Error())
			} else {
				for key, value := range list {
					apiList[key] = value
				}
			}
		}
	}
	return apiList, errorList
}
