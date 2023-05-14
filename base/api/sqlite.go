package api

import (
	"encoding/json"
	"fmt"

	apiBase "github.com/crackeer/simple_http"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type SQLiteServiceAPI struct {
	API         string `json:"api"`
	ContentType string `json:"content_type"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Service     string `json:"service"`
	Timeout     int64  `json:"timeout"`
}

// SQLiteService
type SQLiteService struct {
	CodeKey        string `json:"code_key"`
	DataKey        string `json:"data_key"`
	Env            string `json:"env"`
	DisableExtract int    `json:"disable_extract"`
	Host           string `json:"host"`
	MessageKey     string `json:"message_key"`
	Service        string `json:"service"`
	Sign           string `json:"sign"`
	SignConfig     string `json:"sign_config"`
	SuccessCodeKey string `json:"success_code_key"`
	Timeout        int64  `json:"timeout"`
}

// GetServiceAPIFromSQLite
//
//	@param sqliteFile
//	@return map[string]*apiBase.ServiceAPI
//	@return error
func GetServiceAPIFromSQLite(sqliteFile string) (map[string]*apiBase.ServiceAPI, error) {
	sqliteDB, err := gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	serviceList := []SQLiteService{}
	sqliteDB.Model(&SQLiteService{}).Find(&serviceList)

	services := []string{}
	for _, service := range serviceList {
		services = append(services, service.Service)
	}
	apis := []SQLiteServiceAPI{}
	sqliteDB.Model(&SQLiteServiceAPI{}).Where(map[string]interface{}{
		"service": services,
	}).Find(&apis)
	apiGroup := map[string][]SQLiteServiceAPI{}
	for _, item := range apis {
		if _, ok := apiGroup[item.Service]; !ok {
			apiGroup[item.Service] = []SQLiteServiceAPI{}
		}
		apiGroup[item.Service] = append(apiGroup[item.Service], item)
	}
	retData := map[string]*apiBase.ServiceAPI{}
	for _, item := range serviceList {
		if _, ok := apiGroup[item.Service]; !ok {
			continue
		}
		for _, tmpAPI := range apiGroup[item.Service] {
			key := fmt.Sprintf("%s/%s:%s", item.Service, tmpAPI.API, item.Env)
			signConfig := map[string]interface{}{}
			json.Unmarshal([]byte(item.SignConfig), &signConfig)
			retData[key] = &apiBase.ServiceAPI{
				Host:           item.Host,
				DisableExtract: item.DisableExtract > 0,
				SignName:       item.Sign,
				SignConfig:     signConfig,
				Path:           tmpAPI.Path,
				Method:         tmpAPI.Method,
				Timeout:        item.Timeout,
			}
		}
	}
	return retData, nil

}
