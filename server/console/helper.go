package console

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/crackeer/goaway/container"
	"github.com/crackeer/goaway/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/tidwall/gjson"
)

var (
	tokenKey string = "token"
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

// GetCurrentUser
//
//	@param ctx
//	@return model.User
func GetCurrentUser(ctx *gin.Context) model.User {
	value, exist := ctx.Get("CurrentUser")
	if !exist {
		return model.User{}
	}
	if user, ok := value.(model.User); ok {
		return user
	}
	return model.User{}
}

func extractID(data interface{}) int64 {
	bytes, _ := json.Marshal(data)
	return gjson.GetBytes(bytes, "id").Int()
}

// getCookieDomain
//
//	@param ctx
//	@return string
func getCookieDomain(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	host := ctx.Request.Host
	if strings.Contains(host, ":") {
		return strings.Split(host, ":")[0]
	}
	return host
}
