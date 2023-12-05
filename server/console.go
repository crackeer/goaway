package server

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/crackeer/goaway/container"
	"github.com/crackeer/goaway/server/console"
	ginHelper "github.com/crackeer/gopkg/gin"
	apiBase "github.com/crackeer/simple_http"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

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

// RunConsole
func RunConsole() error {
	cfg := container.GetAppConfig()
	router := gin.New()
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false
	router.POST("/user/login", ginHelper.DoResponseJSON(), console.Login)
	wrapperRouter := router.Group("", console.CheckAPILogin, ginHelper.DoResponseJSON())
	wrapperRouter.GET("/user/info", func(ctx *gin.Context) {
		ginHelper.Success(ctx, console.GetCurrentUser(ctx))
	})
	wrapperRouter.POST("/delete/:table/:id", console.Delete, console.RecordLog("delete"))
	wrapperRouter.POST("/create/:table", console.Create, console.RecordLog("create"))
	wrapperRouter.POST("/modify/:table/:id", console.Modify, console.RecordLog("modify"))
	wrapperRouter.GET("/query/:table", console.Query)
	wrapperRouter.GET("/env/list", func(ctx *gin.Context) {
		ginHelper.Success(ctx, container.GetAppConfig().EnvList)
	})
	wrapperRouter.GET("/router/category", func(ctx *gin.Context) {
		ginHelper.Success(ctx, container.GetAppConfig().RouterCategory)
	})
	wrapperRouter.GET("/sign/list", getSignList)
	router.Use(console.CheckLogin)
	router.NoRoute(createStaticHandler(http.Dir(cfg.StaticDir)))
	return router.Run(fmt.Sprintf(":%d", cfg.ConsolePort))
}

func getSignList(ctx *gin.Context) {
	list := apiBase.GetSignHandleMap()
	retData := []map[string]interface{}{}
	for _, v := range list {
		retData = append(retData, map[string]interface{}{
			"sign_id":         v.ID(),
			"introduction":    v.Introduction(),
			"config_template": v.SignConfigTemplate(),
			"go_pkg_path":     reflect.TypeOf(v).PkgPath(),
			"name":            reflect.TypeOf(v).Name(),
		})
	}
	ginHelper.Success(ctx, retData)
}
