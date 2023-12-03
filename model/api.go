package model

import (
	"encoding/json"
	"fmt"

	apiBase "github.com/crackeer/simple_http"
	"gorm.io/gorm"
)

type ServiceAPI struct {
	ID          int64  `json:"id"`
	API         string `json:"api"`
	ContentType string `json:"content_type"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Service     string `json:"service"`
	Description string `json:"description"`
	Timeout     int64  `json:"timeout"`
}

func (ServiceAPI) TableName() string {
	return "service_api"
}

// Service
type Service struct {
	ID             int64  `json:"id"`
	CodeKey        string `json:"code_key"`
	DataKey        string `json:"data_key"`
	Env            string `json:"env"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	DisableExtract int    `json:"disable_extract"`
	Host           string `json:"host"`
	MessageKey     string `json:"message_key"`
	Service        string `json:"service"`
	Sign           string `json:"sign"`
	SignConfig     string `json:"sign_config"`
	SuccessCode    string `json:"success_code"`
	Timeout        int64  `json:"timeout"`
}

func (Service) TableName() string {
	return "service"
}

// GetServiceAPIFrom
//
//	@param File
//	@return map[string]*apiBase.ServiceAPI
//	@return error
func GetServiceAPIFromDB(db *gorm.DB, env string) (map[string]*apiBase.ServiceAPI, error) {
	defaultServiceList := []*Service{}
	envServiceList := []*Service{}
	db.Model(&Service{}).Where(map[string]interface{}{
		"env": "default",
	}).Find(&defaultServiceList)
	if len(env) > 0 {
		db.Model(&Service{}).Where(map[string]interface{}{
			"env": env,
		}).Find(&envServiceList)
	}
	useServiceList := mergeServiceList(defaultServiceList, envServiceList)

	services := []string{}
	for _, service := range useServiceList {
		services = append(services, service.Service)
	}
	apis := []ServiceAPI{}
	db.Model(&ServiceAPI{}).Where(map[string]interface{}{
		"service": services,
	}).Find(&apis)

	apiGroup := map[string][]ServiceAPI{}
	for _, item := range apis {
		if _, ok := apiGroup[item.Service]; !ok {
			apiGroup[item.Service] = []ServiceAPI{}
		}
		apiGroup[item.Service] = append(apiGroup[item.Service], item)
	}
	retData := map[string]*apiBase.ServiceAPI{}
	for _, item := range useServiceList {
		for _, tmpAPI := range apiGroup[item.Service] {
			key := fmt.Sprintf("%s/%s", item.Service, tmpAPI.API)
			signConfig := map[string]interface{}{}
			json.Unmarshal([]byte(item.SignConfig), &signConfig)
			retData[key] = &apiBase.ServiceAPI{
				Host:           item.Host,
				DisableExtract: item.DisableExtract > 0,
				Sign:           item.Sign,
				SignConfig:     signConfig,
				Path:           tmpAPI.Path,
				Method:         tmpAPI.Method,
				Timeout:        item.Timeout,
			}
		}
	}
	return retData, nil
}

func mergeServiceList(defaultServiceList, envServiceList []*Service) []*Service {
	mapService := make(map[string]*Service)
	for _, service := range defaultServiceList {
		mapService[service.Service] = service
	}

	for _, service := range envServiceList {
		mapService[service.Service] = service
	}

	list := []*Service{}
	for _, service := range mapService {
		list = append(list, service)
	}
	return list
}

func init() {
	registerNewModelFunc("service", func() (interface{}, interface{}) {
		return &Service{}, []Service{}
	})
	registerNewModelFunc("service_api", func() (interface{}, interface{}) {
		return &ServiceAPI{}, []ServiceAPI{}
	})
}
