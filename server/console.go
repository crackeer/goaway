package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/crackeer/goaway/container"
	"github.com/crackeer/goaway/model"
	ginHelper "github.com/crackeer/gopkg/gin"
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
	router.POST("/delte/:table/:id", ginHelper.DoResponseJSON(), deleteData)
	router.POST("/create/:table", ginHelper.DoResponseJSON(), deleteData)
	router.POST("/modify/:table/:id", ginHelper.DoResponseJSON(), modifyData)
	router.GET("/query/:table", ginHelper.DoResponseJSON(), queryData)
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
	db := container.GetModelDB()
	result := db.Exec("DELETE FROM ? where id = ?", getTable(ctx), getDataID(ctx))
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
		page      int = 1
		pageSize  int = 20
		list      []map[string]interface{}
		total     int64
		totalPage int = 0
	)
	db := container.GetModelDB()
	query := ginHelper.AllGetParams(ctx)
	if value, ok := strconv.Atoi(query["_page"]); ok == nil && value > 0 {
		page = value
	}
	if value, ok := strconv.Atoi(query["_page_size"]); ok == nil && value > 0 {
		pageSize = value
	}

	delete(query, "_page")
	delete(query, "_page_size")

	offset := (page - 1) * pageSize

	db.Table(getTable(ctx)).Where(query).Order("id desc").Offset(offset).Limit(pageSize).Find(&list)
	db.Table(getTable(ctx)).Where(query).Count(&total)
	totalPage = int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage = totalPage + 1
	}
	ginHelper.Success(ctx, map[string]interface{}{
		"list":         list,
		"current_page": page,
		"page_size":    pageSize,
		"total":        total,
		"total_page":   totalPage,
	})
}
