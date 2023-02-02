package server

import (
	"github.com/crackeer/go-gateway/server/handler"
	"github.com/gin-gonic/gin"
)

func initEndpoint() *gin.Engine {
	router := gin.New()
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false
	setupInternal(router)
	return router
}

// setupInternal
//
//	@param router
func setupInternal(router *gin.Engine) {
	if router == nil {
		return
	}
	router.NoRoute(handler.Handle)
}
