package container

import (
	"fmt"
	"strings"

	"github.com/crackeer/go-gateway/base/sign"
	apiBase "github.com/crackeer/simple_http"
)

// InitSign
//
//	@return error
func InitSign() error {
	data, errorList := getSign(config.SignDir, config.DBConnection)
	if len(errorList) > 0 {
		panic(fmt.Sprintf("get sign error:%s", strings.Join(errorList, ";")))
	}
	for key, value := range data {
		err := apiBase.RegisterLuaSign(key, value)
		if err != nil {
			panic(fmt.Errorf("register sign error: %v[%s]", err.Error(), key))
		}
	}
	return nil
}

func getSign(signDir, dbConnection string) (map[string]string, []string) {
	var (
		retData   map[string]string = map[string]string{}
		errorList []string
	)
	if len(signDir) > 0 {
		tmp, err := sign.GetSignCodeFromDir(signDir)
		if err != nil {
			errorList = append(errorList, err.Error())
		} else {
			for key, value := range tmp {
				retData[key] = value
			}
		}
	}

	if len(dbConnection) > 0 {
		db, err := OpenDatabase(config.DBConnection)
		if err != nil {
			errorList = append(errorList, fmt.Sprintf("connect %s error:%s", config.DBConnection, err.Error()))
		} else {
			tmp, err := sign.GetSignCodeFromDB(db)
			if err != nil {
				panic(err.Error())
			}
			for key, value := range tmp {
				retData[key] = value
			}
		}
	}
	return retData, errorList
}
