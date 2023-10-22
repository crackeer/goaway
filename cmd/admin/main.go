package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/crackeer/goaway/util/database"
	ginHelper "github.com/crackeer/gopkg/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

type AppConfig struct {
	Port      int64  `env:"PORT"`
	Database  string `env:"DATABASE"`
	LogDir    string `env:"LOG_DIR"`
	StaticDir string `env:"STATIC_DIR"`
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

func main() {
	cfg := &AppConfig{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	router := gin.New()
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false
	db, err := database.Open(cfg.Database)
	if err != nil {
		panic(err)
	}
	router.Any("/api/*path", ginHelper.DoResponseJSON(), handleAPI(db))
	router.NoRoute(createStaticHandler(http.Dir(cfg.StaticDir)))
	router.Run(fmt.Sprintf(":%d", cfg.Port))
}
