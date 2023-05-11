package container

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	api "github.com/crackeer/simple_http"
)

type APIContent struct {
	Service []struct {
		Env        string                 `json:"env"`
		Host       string                 `json:"host"`
		Timeout    int64                  `json:"timeout"`
		SignConfig map[string]interface{} `json:"sign_config"`
	} `json:"service"`
	DisableExtract bool   `json:"disable_extract"`
	DataKey        string `json:"data_key"`
	CodeKey        string `json:"code_key"`
	SuccessCode    string `json:"success_code"`
	MessageKey     string `json:"message_key"`
	SignName       string `json:"sign_name"`
	APIList        []struct {
		Name        string `json:"name"`
		Path        string `json:"path"`
		Method      string `json:"method"`
		ContentType string `json:"content_type"`
	} `json:"api_list"`
}

func registerAPI(dir string) error {
	list, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, name := range list {
		content, err := getAPIContent(filepath.Join(dir, name.Name()))
		if err != nil {
			return fmt.Errorf("read api content error: %v[%s]", err, filepath.Join(dir, name.Name()))
		}
		parts := strings.Split(name.Name(), ".")
		registerAPIContent(content, parts[0])
	}
	return nil
}

func getAPIContent(file string) (*APIContent, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	apiContent := &APIContent{}
	if err := json.Unmarshal(bytes, apiContent); err != nil {
		return nil, err
	}
	return apiContent, nil
}

func registerAPIContent(content *APIContent, name string) error {
	list := map[string]*api.ServiceAPI{}
	for _, service := range content.Service {
		for _, item := range content.APIList {
			key := fmt.Sprintf("%s/%s/%s", service.Env, name, item.Name)
			fmt.Println(key)
			list[key] = &api.ServiceAPI{
				Host:           service.Host,
				DisableExtract: content.DisableExtract,
				SignName:       content.SignName,
				SignConfig:     service.SignConfig,
				Path:           item.Path,
				Method:         item.Method,
				Timeout:        service.Timeout,
			}
		}
	}
	for key, value := range list {
		api.RegisterServiceAPI(key, value)
	}
	return nil
}

func registerSign(dir string) error {
	list, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, name := range list {
		parts := strings.Split(name.Name(), ".")
		err := api.RegisterLuaSignByFile(parts[0], filepath.Join(dir, name.Name()))
		if err != nil {
			return fmt.Errorf("register sign error: %v[%s]", err, filepath.Join(dir, name.Name()))
		}
	}
	return nil
}

// InitAPI
//
//	@param cfg
//	@return error
func InitAPI(cfg *AppConfig) error {
	if err := registerAPI(cfg.APIDir); err != nil {
		panic(err.Error())
	}
	if err := registerSign(cfg.SignDir); err != nil {
		panic(err.Error())
	}
	/*
			api.RegisterServiceAPI("abc/test", &api.ServiceAPI{
				Host:           "https://www.boredapi.com",
				DisableExtract: true,
				SignName:       "test",
				SignConfig: map[string]interface{}{
					"ak": "22",
				},
				Path:    "api/activity",
				Method:  http.MethodGet,
				Timeout: 3000,
			})
		err := api.RegisterLuaSignByFile("test", "./config/sign/test.lua")
		if err != nil {
			panic(err.Error())
		}*/
	return nil
}
