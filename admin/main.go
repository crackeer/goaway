package admin

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/crackeer/go-gateway/util/database"
	ginHelper "github.com/crackeer/gopkg/gin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminConfig
type AdminConfig struct {
	Port      int64
	DB        *gorm.DB
	StaticDir string
}

// Run
//
//	@param ctx
func Run(ctx context.Context, cfg *AdminConfig) error {
	router := gin.New()
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false
	router.Any("/api/*path", ginHelper.DoResponseJSON(), handleAPI(cfg.DB))
	router.NoRoute(createStaticHandler(http.Dir(cfg.StaticDir)))
	return router.Run(fmt.Sprintf(":%d", cfg.Port))
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

func handleAPI(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req, err := database.ParseRequest(ctx.Request, ctx.Param("path"))
		if err != nil {
			ginHelper.Failure(ctx, -1, err.Error())
			return
		}
		req.UseDB(db, "mysql")
		result, err := req.Handle()
		if err != nil {
			ginHelper.Failure(ctx, -1, err.Error())
			return
		}
		ginHelper.Success(ctx, result)
	}
}
