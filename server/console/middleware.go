package console

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/crackeer/goaway/container"
	"github.com/crackeer/goaway/model"
	ginHelper "github.com/crackeer/gopkg/gin"
	"github.com/gin-gonic/gin"
)

// RecordLog
//
//	@param action
//	@return gin.HandlerFunc
func RecordLog(action string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		table := getTable(ctx)
		user := GetCurrentUser(ctx)
		log := model.Log{
			Action:   action,
			Table:    table,
			CreateAt: time.Now().Unix(),
			ModifyAt: time.Now().Unix(),
			UserID:   user.ID,
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

// CheckLogin
//
//	@param ctx
func CheckLogin(ctx *gin.Context) {
	if !strings.HasSuffix(ctx.Request.URL.Path, ".html") || strings.HasSuffix(ctx.Request.URL.Path, "user/login.html") {
		return
	}

	redirectLogin := func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "/user/login.html?jump="+ctx.Request.URL.Path)
		ctx.Abort()
	}

	token, err := ctx.Cookie(tokenKey)
	if err != nil {
		redirectLogin(ctx)
		return
	}
	loginUser, err := parseJwt(token)
	if err != nil {
		redirectLogin(ctx)
		return
	}
	if loginUser.ExpiresAt > time.Now().Unix() {
		redirectLogin(ctx)
		return
	}
}

// CheckAPILogin
//
//	@param ctx
func CheckAPILogin(ctx *gin.Context) {

	token, err := ctx.Cookie(tokenKey)
	if err != nil {
		ginHelper.Failure(ctx, -100, "user not login")
		return
	}
	loginUser, err := parseJwt(token)
	if err != nil {
		ginHelper.Failure(ctx, -100, "user not login")
		return
	}
	if loginUser.ExpiresAt > time.Now().Unix() {
		ginHelper.Failure(ctx, -100, "user not login")
		return
	}
	ctx.Set("CurrentUser", loginUser.User)
}

// CheckPermission
//
//	@param ctx
func CheckPermission(action string) gin.HandlerFunc {
	mapping := map[string]map[string]string{
		"router": map[string]string{
			"create": "writer",
			"modify": "writer",
			"query":  "writer,reader",
		},
		"service": map[string]string{
			"create": "writer",
			"modify": "writer",
			"query":  "writer,reader",
		},
		"service_api": map[string]string{
			"create": "writer",
			"modify": "writer",
			"query":  "writer,reader",
		},
	}
	return func(ctx *gin.Context) {
		user := GetCurrentUser(ctx)
		if user.UserType == model.UserTypeRoot {
			return
		}

		permission, ok := mapping[ctx.Param("table")]; !ok {
			ginHelper.Failure(ctx, -90, "user not allowed")
			return
		}

		userTypes, ok := permission[acaction]; !ok {
			ginHelper.Failure(ctx, -90, "user not allowed")
			return
		}
		var allowed bool
		for _, item ：= strings.Split(userTypes, ",") {
			if item == user.UserType {
				allowed = true
				break
			}
		}

		if !allowed {
			ginHelper.Failure(ctx, -90, "user not allowed")
			return
		}
	}
}
