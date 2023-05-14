package handler

import (
	"fmt"
	"strings"

	"github.com/crackeer/go-gateway/container"
	giner "github.com/crackeer/gopkg/gin"
	api "github.com/crackeer/simple_http"
	"github.com/gin-gonic/gin"
)

// Handle 处理
//
//	@param ctx
func Handle(ctx *gin.Context) {
	path := ctx.Request.URL.String()
	path = strings.Trim(path, "/")
	config, err := container.GetRouter(path)
	if err != nil {
		giner.Failure(ctx, -1, err.Error())
		return
	}
	if config.Mode == "static" {
		giner.Success(ctx, config.Response)
		return
	}

	input := giner.AllParams(ctx)
	header := giner.AllHeader(ctx)
	fmt.Println()
	response, err := api.RequestServiceAPIByName(config.ProxyAPI, input, header)
	if err != nil {
		giner.Failure(ctx, -1, err.Error())
		return
	}

	if response.Error {
		giner.Failure(ctx, -1, response.Message)
		return
	}

	giner.Success(ctx, response.Data)
}
