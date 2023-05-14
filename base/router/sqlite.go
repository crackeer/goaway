package router

import (
	"encoding/json"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// SqliteRouter
type SqliteRouter struct {
	Mode     string `json:"mode"`
	Path     string `json:"path"`
	Request  string `json:"request"`
	Response string `json:"response"`
	Status   int    `json:"status"`
}

// GetRouterFromSQLite
//
//	@param sqliteFile
//	@return map
func GetRouterFromSQLite(sqliteFile string) (map[string]*RouterConfig, error) {
	sqliteDB, err := gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	list := []SqliteRouter{}
	sqliteDB.Model(&SqliteRouter{}).Find(&list)
	retData := map[string]*RouterConfig{}
	for _, item := range list {
		path := strings.Trim(item.Path, "/")
		tmp := &RouterConfig{
			Mode:     item.Mode,
			ProxyAPI: item.Request,
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
