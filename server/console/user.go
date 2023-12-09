package console

import (
	"net/http"
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

func Logout(ctx *gin.Context) {
	domain := getCookieDomain(ctx)
	ctx.SetCookie(tokenKey, "1", -1, "/", domain, true, false)
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

// Register
//
//	@param ctx
func Register(ctx *gin.Context) {
	currentUser := GetCurrentUser(ctx)
	if currentUser.UserType != model.UserTypeRoot {
		ginHelper.Failure(ctx, -1, "暂无权限创建")
		return
	}
	registerUser := &UserRegister{}
	if err := ctx.ShouldBindJSON(registerUser); err != nil {
		ginHelper.Failure(ctx, -2, err.Error())
		return
	}
	db := container.GetModelDB()
	tmpUser := model.User{}
	if err := db.Model(&model.User{}).Where("username = ?", registerUser.Username).First(&tmpUser).Error; err == nil && tmpUser.ID > 0 {
		ginHelper.Failure(ctx, -3, "用户已存在")
		return
	}

	user := &model.User{
		Username:    registerUser.Username,
		UserType:    registerUser.UserType,
		PasswordMD5: calcMD5(registerUser.Password),
	}

	if err := db.Create(&user).Error; err != nil {
		ginHelper.Failure(ctx, -2, err.Error())
		return
	}

	ginHelper.Success(ctx, map[string]interface{}{
		"username":  registerUser.Username,
		"user_type": registerUser.UserType,
	})
}

type UserInfo struct {
	Username        string `json:"username"`
	Role            string `json:"role"`
	RoleName        string `json:"role_name"`
	Superuser       bool   `json:"superuser"`
	RoleDescription string `json:"role_description"`
}

// GetUserInfo
//
//	@param ctx
func GetUserInfo(ctx *gin.Context) {
	user := GetCurrentUser(ctx)
	permission := container.GetPermission()
	userInfo := map[string]interface{}{

		"username":         user.Username,
		"role":             user.UserType,
		"created_at":       user.CreateAt,
		"role_name":        "暂无",
		"superuser":        false,
		"role_description": "",
	}
	for _, role := range permission.Roles {
		if role.Role == user.UserType {
			userInfo["role_name"] = role.Name
			userInfo["superuser"] = role.Superuser
			userInfo["role_description"] = role.Description
			break
		}
	}
	userPermission := map[string]string{}
	for key, tables := range permission.Permissions {
		parts := strings.Split(key, ":")
		if parts[0] == user.UserType && len(parts) >= 2 {
			userPermission[parts[1]] = tables
		}
	}
	userInfo["permission"] = userPermission
	ginHelper.Success(ctx, userInfo)
}

// FormateUser
//
//	@param user
//	@return *UserInfo
func FormateUser(user *model.User) *UserInfo {
	retData := &UserInfo{
		Username: user.Username,
		Role:     user.UserType,
	}
	return retData
}
