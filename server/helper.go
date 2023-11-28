package server

import (
	"fmt"

	"github.com/crackeer/goaway/container"
	"github.com/crackeer/goaway/model"
	"github.com/golang-jwt/jwt"
)

// generateJwt
//
//	@param user
//	@return string
//	@return error
func generateJwt(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(user.Map()))
	return token.SignedString([]byte(container.GetAppConfig().LoginSalt))
}

// parseJwt
//
//	@param token
//	@return map[string]interface{}
//	@return error
func parseJwt(token string) (model.User, error) {
	object, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(container.GetAppConfig().LoginSalt), nil
	})

	if err != nil {
		return model.User{}, err
	}
	user := object.Claims.(jwt.MapClaims)
	return model.Map2User(user), nil
}
