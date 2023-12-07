package model

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

type NewModelFunc func() (interface{}, interface{})

var newModelFuncMap map[string]NewModelFunc = map[string]NewModelFunc{}

var locker *sync.Mutex = &sync.Mutex{}

// registerNewModelFunc
//
//	@param table
//	@param someFunc
func registerNewModelFunc(table string, someFunc NewModelFunc) {
	locker.Lock()
	defer locker.Unlock()
	newModelFuncMap[table] = someFunc
}

// NewModel
//
//	@param table
//	@return interface{}
func NewModel(table string) (interface{}, interface{}) {
	if someFunc, ok := newModelFuncMap[table]; ok {
		return someFunc()
	}
	return nil, nil
}

// CheckModelCreate
//
//	@param db
//	@param value
//	@return error
func CheckModelCreate(db *gorm.DB, value interface{}) error {
	if router, ok := value.(*Router); ok {
		return checkRouterCreate(db, router)
	}
	if service, ok := value.(*Service); ok {
		return checkServiceCreate(db, service)
	}
	if api, ok := value.(*ServiceAPI); ok {
		return checkServiceAPICreate(db, api)
	}
	return fmt.Errorf("not a valid model")
}

// MakeModifyData
//
//	@param table
//	@param data
//	@return map[string]interface{}
//	@return nil
func MakeModifyData(table string, data map[string]interface{}) (map[string]interface{}, error) {

	switch table {
	case "router":
		delete(data, "path")
	case "service":
		delete(data, "service")
		delete(data, "env")
	case "service_api":
		delete(data, "service")
	default:
		return nil, fmt.Errorf("not valid model")
	}
	return data, nil
}

// CheckModelModify
//
//	@param db
//	@param value
//	@return error
func CheckModelModify(db *gorm.DB, value interface{}, id int64) error {
	if router, ok := value.(*Router); ok {
		router.ID = id
		return checkRouterModify(db, router)
	}
	if service, ok := value.(*Service); ok {
		service.ID = id
		return checkServiceModify(db, service)
	}
	if api, ok := value.(*ServiceAPI); ok {
		api.ID = id
		return checkServiceAPIModify(db, api)
	}
	return fmt.Errorf("not a valid model")
}
