package model

import (
	"gorm.io/gorm"
)

// Sign SQLiteSign
type Sign struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (Sign) TableName() string {
	return "sign"
}

// GetSignCodeFromSQLite GetRouterFromSQLite
//
//	@param sqliteFile
//	@return map[string]*RouterConfig
//	@return error
func GetSignCodeFromDB(db *gorm.DB) (map[string]string, error) {
	list := []Sign{}
	db.Model(&Sign{}).Find(&list)
	retData := map[string]string{}
	for _, item := range list {
		retData[item.Name] = item.Code
	}

	return retData, nil
}
