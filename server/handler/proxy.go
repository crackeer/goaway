package handler

import (
	"strings"

	giner "github.com/crackeer/gopkg/gin"
	api "github.com/crackeer/simple_http"
	"github.com/gin-gonic/gin"
)

// Handle 处理
//
//	@param ctx
func Proxy(ctx *gin.Context) {
	apiName := ctx.Param("api")
	apiName = strings.TrimLeft(apiName, "/")
	input := giner.AllParams(ctx)
	header := giner.AllHeader(ctx)
	response, err := api.RequestByName(apiName, input, header)
	if err != nil {
		giner.Success(ctx, err.Error())
	} else {
		giner.Success(ctx, response)
	}

}
