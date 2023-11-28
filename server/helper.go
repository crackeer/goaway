package server

import (
	"fmt"

	"github.com/crackeer/goaway/container"
	"github.com/golang-jwt/jwt"
)

// generateJwt
//
//	@param user
//	@return string
//	@return error
func generateJwt(user map[string]interface{}) (map[string]interface{}, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(user))
	return token.SignedString([]byte(container.GetAppConfig().LoginSalt))
}

func parseJwt(token string) (map[string]interface{}, error) {
	object, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(container.GetAppConfig().LoginSalt), nil
	})

	if err != nil {
		return nil, err
	}
	return object.Claims.(jwt.MapClaims), nil
}
