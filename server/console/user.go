package console

import (
	"strings"
	"time"

	"github.com/crackeer/goaway/container"
	"github.com/crackeer/goaway/model"
	ginHelper "github.com/crackeer/gopkg/gin"
	"github.com/gin-gonic/gin"
)

// LoginForm
type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserRegister
type UserRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

// Login
//
//	@param ctx
func Login(ctx *gin.Context) {
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

// Register
//
//	@param ctx
func Register(ctx *gin.Context) {
	currentUser := GetCurrentUser(ctx)
	if currentUser.UserType != model.UserTypeRoot {
		ginHelper.Failure(ctx, -1, "not allowed to register user")
		return
	}
	registerUser := &UserRegister{}
	if err := ctx.ShouldBindJSON(registerUser); err != nil {
		ginHelper.Failure(ctx, -1, err.Error())
		return
	}
	user := &model.User{
		Username:    registerUser.Username,
		UserType:    registerUser.UserType,
		PasswordMD5: calcMD5(registerUser.Password),
	}
	db := container.GetModelDB()
	if err := db.Create(&user).Error; err != nil {
		ginHelper.Failure(ctx, -2, err.Error())
		return
	}

	ginHelper.Success(ctx, map[string]interface{}{
		"username":  registerUser.Username,
		"user_type": registerUser.UserType,
	})
}
