package handler

import (
	"time"

	"github.com/crackeer/go-gateway/container"
	"github.com/crackeer/go-gateway/util"
	ginHelper "github.com/crackeer/gopkg/gin"
	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	AppID  string `json:"app_id"`
	Secret string `json:"secret"`
}

// Token
//
//	@param ctx
func Token(ctx *gin.Context) {
	request := &TokenRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ginHelper.Failure(ctx, -1, err.Error())
		ctx.Abort()
		return
	}
	app, err := container.GetApp(request.AppID)
	if err != nil {
		ginHelper.Failure(ctx, -1, err.Error())
		return
	}
	if app.Secret != request.Secret {
		ginHelper.Failure(ctx, -1, "secret not right")
		return
	}

	expire := time.Now().Add(300 * time.Second).Unix()
	token, err := util.GenerateToken(app.Secret, &util.PlainObject{
		AppID:    request.AppID,
		ExpireAt: expire,
	})
	if err != nil {
		ginHelper.Failure(ctx, -1, err.Error())
		return

	}

	ginHelper.Success(ctx, map[string]interface{}{
		"token":     token,
		"expire_at": expire,
	})
}
