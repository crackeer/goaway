package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminConfig
type AdminConfig struct {
	Port      int64
	DSN       string
	StaticDir string
}

// Run
//
//	@param ctx
func Run(ctx context.Context, cfg *AdminConfig) error {
	router := gin.New()
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false
	router.NoRoute(createStaticHandler(http.Dir(cfg.StaticDir)))
}

func createStaticHandler(fs http.FileSystem) gin.HandlerFunc {
	fileServer := http.StripPrefix("", http.FileServer(fs))
	return func(ctx *gin.Context) {
		file := strings.TrimLeft(ctx.Request.URL.Path, "/")
		f, err := fs.Open(file)
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusNotFound)
			ctx.Abort()
			return
		}
		f.Close()
		fileServer.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
