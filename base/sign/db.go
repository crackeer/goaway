package sign

import (
	"gorm.io/gorm"
)

// SQLiteSign
type SQLiteSign struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (SQLiteSign) TableName() string {
	return "sign"
}

// GetSignCodeFromSQLite GetRouterFromSQLite
//
//	@param sqliteFile
//	@return map[string]*RouterConfig
//	@return error
func GetSignCodeFromDB(db *gorm.DB) (map[string]string, error) {
	list := []SQLiteSign{}
	db.Model(&SQLiteSign{}).Find(&list)
	retData := map[string]string{}
	for _, item := range list {
		retData[item.Name] = item.Code
	}

	return retData, nil
}
