package handler

import (
	giner "github.com/crackeer/gopkg/gin"
	"github.com/gin-gonic/gin"
)

// Handle 处理
//
//	@param ctx
func Handle(ctx *gin.Context) {
	giner.Success(ctx, map[string]interface{}{})
}
