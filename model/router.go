package model

import (
	"encoding/json"
	"strings"

	"gorm.io/gorm"
)

// RouterConfig
type RouterConfig struct {
	Mode     string            `json:"mode"`
	ProxyAPI string            `json:"proxy_api"`
	Input    interface{}       `json:"input"`
	Header   map[string]string `json:"header"`
	Response interface{}       `json:"response"`
	Status   int               `json:"status"`
}

// SqliteRouter
type Router struct {
	ID          int64  `json:"id"`
	Mode        string `json:"mode"`
	Category    string `json:"category"`
	Path        string `json:"path"`
	ServiceAPI  string `json:"service_api"`
	Input       string `json:"input"`
	Header      string `json:"header"`
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
			ProxyAPI: item.ServiceAPI,
			Status:   item.Status,
		}
		var (
			response interface{}
			input    interface{}
			header   map[string]string = make(map[string]string)
		)
		if err := json.Unmarshal([]byte(item.Response), &response); err == nil {
			tmp.Response = response
		} else {
			tmp.Response = item.Response
		}
		if err := json.Unmarshal([]byte(item.Input), &input); err == nil {
			tmp.Input = input
		}
		if err := json.Unmarshal([]byte(item.Header), &header); err == nil {
			tmp.Header = header
		}
		retData[path] = tmp
	}

	return retData, nil
}

func init() {
	registerNewModelFunc("router", func() (interface{}, interface{}) {
		return &Router{}, []Router{}
	})
}
