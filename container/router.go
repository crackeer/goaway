package container

import (
	"fmt"

	"github.com/crackeer/gopkg/router"
	"github.com/crackeer/gopkg/router/api"
)

var (
	routerFactory router.RouterFactory
	apiFactory    api.APIFactory
)

// InitRouterFactory
//
//	@param cfg
func InitRouterFactory(cfg *AppConfig) error {
	factory1, err := router.NewFileRouter(cfg.RouterDir)
	if err != nil {
		return fmt.Errorf("init router factory error: %v", err.Error())
	}
	routerFactory = factory1
	factory2, _ := api.NewJSONAPI(cfg.APIDir)
	apiFactory = factory2
	return nil
}

func GetAPIFacory() api.APIFactory {
	return apiFactory
}

func GetRouterFactory() router.RouterFactory {
	return routerFactory
}
