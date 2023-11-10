package model

import (
	"encoding/json"
	"strings"

	"gorm.io/gorm"
)

// RouterConfig
type RouterConfig struct {
	Mode     string      `json:"mode"`
	ProxyAPI string      `json:"proxy_api"`
	Response interface{} `json:"response"`
	Status   int         `json:"status"`
}

// SqliteRouter
type Router struct {
	Mode        string `json:"mode"`
	Category    string `json:"category"`
	Path        string `json:"path"`
	Request     string `json:"request"`
	Response    string `json:"response"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}

func (Router) TableName() string {
	return "router"
}

// GetRouterFromSQLite
//
//	@param sqliteFile
//	@return map
func GetRouterFromDB(db *gorm.DB) (map[string]*RouterConfig, error) {
	list := []Router{}
	db.Model(&Router{}).Find(&list)
	retData := map[string]*RouterConfig{}
	for _, item := range list {
		path := strings.Trim(item.Path, "/")
		tmp := &RouterConfig{
			Mode:     item.Mode,
			ProxyAPI: item.Request,
			Status:   item.Status,
		}
		var response interface{}
		if err := json.Unmarshal([]byte(item.Response), &response); err == nil {
			tmp.Response = response
		} else {
			tmp.Response = item.Response
		}
		retData[path] = tmp
	}

	return retData, nil
}
