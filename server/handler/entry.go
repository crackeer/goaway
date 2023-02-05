package handler

import (
	"fmt"
	"strings"

	"github.com/crackeer/go-gateway/container"
	giner "github.com/crackeer/gopkg/gin"
	"github.com/crackeer/gopkg/router"
	"github.com/crackeer/gopkg/router/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// NewRouterHander
//
//	@param env
//	@param routerFactory
//	@param apiFactory
//	@return gin.HandlerFunc
func NewRouterHander(env string, routerFactory router.RouterFactory, apiFactory api.APIFactory, logger *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uriPath := strings.TrimLeft(ctx.Request.URL.Path, "/")
		meta := routerFactory.Get(uriPath)
		if meta == nil {
			return
		}
		query := giner.AllParams(ctx)
		executor := router.NewRouterExecuter(apiFactory).UseEnv(env).UseInput(query).UseLogrusLogger(logger)
		if err := executor.Exec(meta); err != nil {
			giner.Failure(ctx, -1, err.Error())
			return
		}

		retData := executor.BuildResponse(meta)

		giner.Success(ctx, retData)
	}
}

// Handle 处理
//
//	@param ctx
func Handle(ctx *gin.Context) {
	appConfig := container.GetAppConfig()
	routerFactory := container.GetRouterFactory()
	uriPath := strings.TrimLeft(ctx.Request.URL.Path, "/")
	fmt.Println(uriPath)
	meta := routerFactory.Get(uriPath)
	if meta == nil {
		return
	}
	query := giner.AllParams(ctx)
	executor := router.NewRouterExecuter(container.GetAPIFacory()).UseEnv(appConfig.Env).UseInput(query)
	if err := executor.Exec(meta); err != nil {
		giner.Failure(ctx, -1, err.Error())
		return
	}

	retData := executor.BuildResponse(meta)

	giner.Success(ctx, retData)
}
