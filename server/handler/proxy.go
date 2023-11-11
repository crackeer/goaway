package handler

import (
	"strings"

	ginHelper "github.com/crackeer/gopkg/gin"
	api "github.com/crackeer/simple_http"
	"github.com/gin-gonic/gin"
)

// Handle 处理
//
//	@param ctx
func Proxy(ctx *gin.Context) {
	apiName := ctx.Param("api")
	apiName = strings.TrimLeft(apiName, "/")
	input := ginHelper.AllParams(ctx)
	header := ginHelper.AllHeader(ctx)
	response, err := api.RequestByName(apiName, input, header)
	if err != nil {
		ginHelper.Success(ctx, err.Error())
	} else {
		ginHelper.Success(ctx, response)
	}

}
