package model

import "sync"

type NewModelFunc func() (interface{}, interface{})

var newModelFuncMap map[string]NewModelFunc = map[string]NewModelFunc{}

var locker *sync.Mutex = &sync.Mutex{}

// registerNewModelFunc
//
//  @param table
//  @param someFunc
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
