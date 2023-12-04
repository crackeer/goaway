package gateway

import (
	"encoding/json"
	"strings"

	"github.com/crackeer/goaway/container"
	ginHelper "github.com/crackeer/gopkg/gin"
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
		ginHelper.Failure(ctx, -1, err.Error())
		return
	}
	if config.Mode == "static" {
		ginHelper.Success(ctx, config.Response)
		return
	}

	input := ginHelper.AllParams(ctx)
	header := ginHelper.AllHeader(ctx)
	data, _ := json.Marshal(map[string]interface{}{
		"input":  input,
		"header": header,
	})
	if config.Header != nil {
		tmp := jsonValue(config.Input, data)
		if value, ok := tmp.(map[string]string); ok {
			header = value
		}
	}
	if config.Input != nil {
		tmp := jsonValue(config.Input, data)
		if value, ok := tmp.(map[string]interface{}); ok {
			input = value
		}
	}

	response, err := api.RequestByName(config.ProxyAPI, input, header)
	if err != nil {
		ginHelper.Failure(ctx, -1, err.Error())
		return
	}

	if response.Error {
		ginHelper.Failure(ctx, -1, response.Message)
		return
	}

	ginHelper.Success(ctx, response.Data)
}
