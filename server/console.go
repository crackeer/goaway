package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/crackeer/goaway/container"
	"github.com/crackeer/goaway/model"
	ginHelper "github.com/crackeer/gopkg/gin"
	apiBase "github.com/crackeer/simple_http"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var (
	tokenKey string = "token"
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
	router.POST("/user/login", ginHelper.DoResponseJSON(), userLogin)
	wrapperRouter := router.Group("", checkAPILogin, ginHelper.DoResponseJSON())
	wrapperRouter.GET("/user/info", func(ctx *gin.Context) {
		ginHelper.Success(ctx, getCurrentUser(ctx))
	})
	wrapperRouter.POST("/delete/:table/:id", deleteData, recordLog("delete"))
	wrapperRouter.POST("/create/:table", createData, recordLog("create"))
	wrapperRouter.POST("/modify/:table/:id", modifyData, recordLog("modify"))
	wrapperRouter.GET("/query/:table", queryData)
	wrapperRouter.GET("/env/list", func(ctx *gin.Context) {
		ginHelper.Success(ctx, container.GetAppConfig().EnvList)
	})
	wrapperRouter.GET("/router/category", func(ctx *gin.Context) {
		ginHelper.Success(ctx, container.GetAppConfig().RouterCategory)
	})
	wrapperRouter.GET("/sign/list", getSignList)
	router.Use(checkLogin)
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
		value interface{}
	)
	value, _ = model.NewModel(table)
	if err := ctx.ShouldBindJSON(value); err != nil {
		ginHelper.Failure(ctx, -1, err.Error())
		return
	}
	result := db.Create(value)
	if result.Error != nil {
		ginHelper.Failure(ctx, -1, result.Error.Error())
	} else {
		ctx.Set("data_id", extractID(value))
		ginHelper.Success(ctx, value)
	}
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
			"go_pkg_path":     reflect.TypeOf(v).PkgPath(),
			"name":            reflect.TypeOf(v).Name(),
		})
	}
	ginHelper.Success(ctx, retData)
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func userLogin(ctx *gin.Context) {
	loginForm := &LoginForm{}
	if err := ctx.ShouldBindJSON(loginForm); err != nil {
		ginHelper.Failure(ctx, -1, err.Error())
		return
	}
	user := &model.User{}
	db := container.GetModelDB()
	db.Model(&model.User{}).Where(map[string]interface{}{
		"username": loginForm.Username,
	}).First(user)
	if user.ID < 1 {
		ginHelper.Failure(ctx, -1, "user not found")
		return
	}

	if !strings.EqualFold(calcMD5(loginForm.Password), user.PasswordMD5) {
		ginHelper.Failure(ctx, -1, "password wrong")
		return
	}

	expireAt := time.Now().Add(30 * 24 * time.Hour).Unix()

	token, err := generateJwt(user, expireAt)
	if err != nil {
		ginHelper.Failure(ctx, -1, "generate token error:"+err.Error())
		return
	}
	domain := getCookieDomain(ctx)
	ctx.SetCookie(tokenKey, token, 3600*24*365, "/", domain, true, false)
	ginHelper.Success(ctx, map[string]interface{}{
		"token":  token,
		"domain": domain,
	})
}

func recordLog(action string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		table := getTable(ctx)
		log := model.Log{
			Action:   action,
			Table:    table,
			CreateAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		if action == "delete" || action == "modify" {
			log.DataID = getDataID(ctx)
		} else if value, exists := ctx.Get("data_id"); exists {
			log.DataID, _ = value.(int64)
		}
		object, _ := model.NewModel(log.Table)
		db := container.GetModelDB()
		db.Table(log.Table).Where("id = ?", log.DataID).Find(object)
		bytes, _ := json.Marshal(object)
		log.Data = string(bytes)
		db.Create(&log)
	}
}
