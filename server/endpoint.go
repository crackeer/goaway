package server

import (
	"context"
	"fmt"

	"github.com/crackeer/go-gateway/container"
	"github.com/crackeer/go-gateway/server/handler"
	giner "github.com/crackeer/gopkg/gin"
	"github.com/gin-gonic/gin"
)

// Run
//
//	@param ctx
func Run(ctx context.Context, port int64) error {
	engine := newGinEngine()
	addRouterHandler(engine)
	return engine.Run(fmt.Sprintf(":%d", port))
}

func newGinEngine() *gin.Engine {
	router := gin.New()
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false
	return router
}

// addRouterHandler ...
//
//	@param router
func addRouterHandler(router *gin.Engine) {
	if router == nil {
		return
	}
	router.Use(giner.DoResponseJSON())
	appConfig := container.GetAppConfig()
	routerFactory := container.GetRouterFactory()
	router.NoRoute(handler.NewRouterHander(appConfig.Env, routerFactory, container.GetAPIFacory(), container.GetLogger(container.LogTypeAPI)))
}
