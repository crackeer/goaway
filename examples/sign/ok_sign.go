package sign

import (
	apiBase "github.com/crackeer/simple_http"
)

func init() {
	apiBase.RegisterSign(OkSign{})
}

type OkSign struct {
}

func (OkSign) ID() string {
	return "ok"
}

func (OkSign) Introduction() string {
	return "ok sign introduction"
}

func (OkSign) SignConfigTemplate() map[string]interface{} {
	return map[string]interface{}{
		"ak": "xxx",
		"sk": "xxxxxx",
	}
}

func (OkSign) Sign(api *apiBase.ServiceAPI, input map[string]interface{}, header map[string]string) (*apiBase.ServiceAPI, map[string]interface{}, map[string]string, error) {
	return api, input, header, nil
}
