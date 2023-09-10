package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// App
type App struct {
	CreateAt    time.Time `json:"create_at"`
	Description string    `json:"description"`
	ID          int       `json:"id"`
	ModifyAt    time.Time `json:"modify_at"`
	Name        string    `json:"name"`
	Secret      string    `json:"secret"`
	Status      int       `json:"status"`
}

// AppConfig
type AppConfig struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

func (App) TableName() string {
	return "app"
}

// GetAppList
//
//	@param db
//	@return map[string]*AppConfig
//	@return error
func GetAppList(db *gorm.DB) (map[string]*AppConfig, error) {
	list := []App{}
	db.Model(&App{}).Find(&list)
	retData := map[string]*AppConfig{}
	for _, item := range list {
		tmp := &AppConfig{
			Name:   item.Name,
			Secret: item.Secret,
		}

		retData[fmt.Sprintf("%d", item.ID)] = tmp
	}

	return retData, nil
}
