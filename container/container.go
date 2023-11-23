package container

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	middlewares  []gin.HandlerFunc
	nakedRouters map[string]gin.HandlerFunc
	routers      map[string]gin.HandlerFunc
	locker       *sync.Mutex
)

func init() {
	middlewares = []gin.HandlerFunc{}
	locker = &sync.Mutex{}
	nakedRouters = map[string]gin.HandlerFunc{}
	routers = map[string]gin.HandlerFunc{}
}

// GetMiddlewares
//
//	@return []gin.HandlerFunc
func GetMiddlewares() []gin.HandlerFunc {
	return middlewares
}

// GetNakedRouters
//
//	@return map
func GetNakedRouters() map[string]gin.HandlerFunc {
	return nakedRouters
}

// GetRouters
//
//	@return map
func GetRouters() map[string]gin.HandlerFunc {
	return routers
}

// RegisterNakedMiddleware
//
//	@param middleware
//	@return error
func RegisterMiddleware(middleware gin.HandlerFunc) {
	locker.Lock()
	defer locker.Unlock()
	middlewares = append(middlewares, middleware)
}

// RegisterNakedRouter
//
//	@param path
//	@param dataFunc
func RegisterNakedRouter(path string, dataFunc gin.HandlerFunc) {
	locker.Lock()
	defer locker.Unlock()
	nakedRouters[path] = dataFunc
}

// RegisterRouter
//
//	@param path
//	@param dataFunc
func RegisterRouter(path string, dataFunc gin.HandlerFunc) {
	locker.Lock()
	defer locker.Unlock()
	routers[path] = dataFunc
}
