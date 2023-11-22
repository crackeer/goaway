package sign

import (
	apiBase "github.com/crackeer/simple_http"
)

func init() {
	apiBase.RegisterSign(SimpleSign{})
}

type SimpleSign struct {
}

func (SimpleSign) ID() string {
	return "simple_sign"
}

func (SimpleSign) Introduction() string {
	return "## simple sign introduction"
}

func (SimpleSign) SignConfigTemplate() map[string]interface{} {
	return map[string]interface{}{
		"sss":    "xxx",
		"ssssss": "xxxxxx",
	}
}

func (SimpleSign) Sign(api *apiBase.ServiceAPI, input map[string]interface{}, header map[string]string) (*apiBase.ServiceAPI, map[string]interface{}, map[string]string, error) {
	return api, input, header, nil
}
