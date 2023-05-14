package router

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/crackeer/gopkg/util"
)

// RouterConfig
type RouterConfig struct {
	Mode     string      `json:"mode"`
	ProxyAPI string      `json:"proxy_api"`
	Response interface{} `json:"response"`
}

// GetRouterConfig
//
//	@param path
//	@return *RouterConfig
//	@return error
func GetRouterFromLocal(routerDir string) (map[string]*RouterConfig, error) {
	list := util.GetFiles(routerDir)
	retData := map[string]*RouterConfig{}
	for _, item := range list {
		c, err := readRouter(item)
		if err != nil {
			return nil, err
		}
		path := strings.Trim("/", strings.TrimLeft(item, routerDir))
		retData[path] = c

	}
	return retData, nil
}

func readRouter(fullPath string) (*RouterConfig, error) {
	bytes, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("read file `%s` error:%s", fullPath, err.Error())
	}
	retData := &RouterConfig{}
	if err := json.Unmarshal(bytes, retData); err != nil {
		return nil, fmt.Errorf("json unmarshal `%s` content error:%s", fullPath, err.Error())
	}
	return retData, nil
}
