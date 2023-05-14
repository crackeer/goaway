package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	apiBase "github.com/crackeer/simple_http"
)

// APIContent
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

// GetServiceAPIFromDir GetAPIFromDir
//
//	@param dir
//	@return map[string]*apiBase.ServiceAPI
//	@return error
func GetServiceAPIFromDir(dir string) (map[string]*apiBase.ServiceAPI, error) {
	list, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	retData := map[string]*apiBase.ServiceAPI{}
	for _, name := range list {
		content, err := getAPIContent(filepath.Join(dir, name.Name()))
		if err != nil {
			return nil, fmt.Errorf("read api content error: %v[%s]", err, filepath.Join(dir, name.Name()))
		}
		parts := strings.Split(name.Name(), ".")
		for key, value := range parseAPIContent(content, parts[0]) {
			retData[key] = value
		}
	}
	return retData, nil
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

func parseAPIContent(content *APIContent, name string) map[string]*apiBase.ServiceAPI {
	list := map[string]*apiBase.ServiceAPI{}
	for _, service := range content.Service {
		for _, item := range content.APIList {
			key := fmt.Sprintf("%s/%s:%s", name, item.Name, service.Env)
			if len(service.Env) < 1 {
				key = fmt.Sprintf("%s/%s", name, item.Name)
			}
			list[key] = &apiBase.ServiceAPI{
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

	return list
}
