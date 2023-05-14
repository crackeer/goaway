package sign

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// SQLiteSign
type SQLiteSign struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// GetSignCodeFromSQLite GetRouterFromSQLite
//
//	@param sqliteFile
//	@return map[string]*RouterConfig
//	@return error
func GetSignCodeFromSQLite(sqliteFile string) (map[string]string, error) {
	sqliteDB, err := gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	list := []SQLiteSign{}
	sqliteDB.Model(&SQLiteSign{}).Find(&list)
	retData := map[string]string{}
	for _, item := range list {
		retData[item.Name] = item.Code
	}

	return retData, nil
}
