package server

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/crackeer/goaway/container"
	"github.com/crackeer/goaway/model"
	ginHelper "github.com/crackeer/gopkg/gin"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// LoginUser
type LoginUser struct {
	model.User
	ExpiresAt int64
}

// generateJwt
//
//	@param user
//	@return string
//	@return error
func generateJwt(user *model.User, expireAt int64) (string, error) {
	data := user.Map()
	data["expires_at"] = expireAt
	delete(data, "password_md5")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(data))
	return token.SignedString([]byte(container.GetAppConfig().LoginSalt))
}

// parseJwt
//
//	@param token
//	@return map[string]interface{}
//	@return error
func parseJwt(token string) (*LoginUser, error) {
	object, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(container.GetAppConfig().LoginSalt), nil
	})

	if err != nil {
		return nil, err
	}
	user := object.Claims.(jwt.MapClaims)
	bytes, _ := json.Marshal(user)
	loginUser := &LoginUser{}
	json.Unmarshal(bytes, loginUser)
	return loginUser, nil
}

func calcMD5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	hashValue := hash.Sum(nil)
	md5Str := hex.EncodeToString(hashValue)
	return md5Str
}

func checkAPILogin(ctx *gin.Context) {

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

func getCurrentUser(ctx *gin.Context) model.User {
	value, exist := ctx.Get("CurrentUser")
	if !exist {
		return model.User{}
	}
	if user, ok := value.(model.User); ok {
		return user
	}
	return model.User{}
}

func checkLogin(ctx *gin.Context) {
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
