package server

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/crackeer/goaway/container"
	"github.com/crackeer/goaway/model"
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
	router.POST("/delete/:table/:id", ginHelper.DoResponseJSON(), deleteData)
	router.POST("/create/:table", ginHelper.DoResponseJSON(), createData)
	router.POST("/modify/:table/:id", ginHelper.DoResponseJSON(), modifyData)
	router.GET("/query/:table", ginHelper.DoResponseJSON(), queryData)
	router.GET("/env/list", ginHelper.DoResponseJSON(), func(ctx *gin.Context) {
		ginHelper.Success(ctx, container.GetAppConfig().EnvList)
	})
	router.GET("/router/category", ginHelper.DoResponseJSON(), func(ctx *gin.Context) {
		ginHelper.Success(ctx, container.GetAppConfig().RouterCategory)
	})
	router.GET("/sign/list", ginHelper.DoResponseJSON(), getSignList)
	router.NoRoute(createStaticHandler(http.Dir(cfg.StaticDir)))
	return router.Run(fmt.Sprintf(":%d", cfg.ConsolePort))
}

func getTable(ctx *gin.Context) string {
	return ctx.Param("table")
}

func getDataID(ctx *gin.Context) int64 {
	id := ctx.Param("id")
	value, _ := strconv.Atoi(id)
	return int64(value)
}

func deleteData(ctx *gin.Context) {
	if dataID := getDataID(ctx); dataID < 1 {
		ginHelper.Failure(ctx, -1, "data id = 0")
		return
	}
	db := container.GetModelDB()
	result := db.Exec(fmt.Sprintf("DELETE FROM %s where id = %d", getTable(ctx), getDataID(ctx)))
	if result.Error != nil {
		ginHelper.Failure(ctx, -1, result.Error.Error())
	} else {
		ginHelper.Success(ctx, map[string]interface{}{
			"affected": result.RowsAffected,
		})
	}
}

func createData(ctx *gin.Context) {
	db := container.GetModelDB()

	var (
		table string = getTable(ctx)
		err   error
		value interface{}
	)
	switch table {
	case "service":
		value, err = bindService(ctx)
	case "service_api":
		value, err = bindServiceAPI(ctx)
	case "router":
		value, err = bindRouter(ctx)
	}
	if err != nil {
		ginHelper.Failure(ctx, -1, err.Error())
		return
	}
	result := db.Create(value)
	if result.Error != nil {
		ginHelper.Failure(ctx, -1, result.Error.Error())
	} else {
		ginHelper.Success(ctx, value)
	}
}

func bindService(ctx *gin.Context) (*model.Service, error) {
	data := &model.Service{}
	if err := ctx.ShouldBindJSON(data); err != nil {
		return nil, err
	}
	return data, nil
}
func bindServiceAPI(ctx *gin.Context) (*model.ServiceAPI, error) {
	data := &model.ServiceAPI{}
	if err := ctx.ShouldBindJSON(data); err != nil {
		return nil, err
	}
	return data, nil
}

func bindRouter(ctx *gin.Context) (*model.Router, error) {
	data := &model.Router{}
	if err := ctx.ShouldBindJSON(data); err != nil {
		return nil, err
	}
	return data, nil
}

func modifyData(ctx *gin.Context) {
	if dataID := getDataID(ctx); dataID < 1 {
		ginHelper.Failure(ctx, -1, "data id = 0")
		return
	}
	db := container.GetModelDB()
	updateData := ginHelper.AllPostParams(ctx)
	result := db.Table(getTable(ctx)).Where(map[string]interface{}{"id": getDataID(ctx)}).Updates(updateData)
	if result.Error != nil {
		ginHelper.Failure(ctx, -1, result.Error.Error())
	} else {
		ginHelper.Success(ctx, map[string]interface{}{
			"affected": result.RowsAffected,
		})
	}
}

func queryData(ctx *gin.Context) {
	var (
		list []map[string]interface{}
	)
	db := container.GetModelDB()
	query := ginHelper.AllGetParams(ctx)

	db.Table(getTable(ctx)).Where(query).Order("id desc").Find(&list)
	ginHelper.Success(ctx, list)
}

func getSignList(ctx *gin.Context) {
	list := apiBase.GetSignHandleMap()
	retData := []map[string]interface{}{}
	for _, v := range list {
		retData = append(retData, map[string]interface{}{
			"sign_id":         v.ID(),
			"introduction":    v.Introduction(),
			"config_template": v.SignConfigTemplate(),
			"from":            reflect.TypeOf(v).PkgPath() + "#" + reflect.TypeOf(v).Name(),
		})
	}
	ginHelper.Success(ctx, retData)
}
