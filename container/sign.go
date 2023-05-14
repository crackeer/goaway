package container

import (
	"fmt"

	"github.com/crackeer/go-gateway/base/sign"
	apiBase "github.com/crackeer/simple_http"
)

// InitSign
//
//	@return error
func InitSign() error {
	data := map[string]string{}
	if len(config.SignDir) > 0 {
		tmp, err := sign.GetSignCodeFromDir(config.SignDir)
		if err != nil {
			panic(err.Error())
		}

		for key, value := range tmp {
			data[key] = value
		}
	}

	if len(config.SqliteFile) > 0 {
		tmp, err := sign.GetSignCodeFromSQLite(config.SqliteFile)
		if err != nil {
			panic(err.Error())
		}

		for key, value := range tmp {
			data[key] = value
		}
	}
	for key, value := range data {
		err := apiBase.RegisterLuaSign(key, value)
		if err != nil {
			panic(fmt.Errorf("register sign error: %v[%s]", err.Error(), key))
		}
	}
	return nil
}
