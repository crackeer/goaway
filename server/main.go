package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/crackeer/goaway/container"
	"github.com/crackeer/goaway/server/gateway"
	giner "github.com/crackeer/gopkg/gin"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
)

func Main() {
	root := context.Background()
	globalWg := &sync.WaitGroup{}
	cancelCtx, cancel := context.WithCancel(root)
	if err := container.InitContainer(cancelCtx, globalWg); err != nil {
		panic(fmt.Sprintf("Failed to initialize container: %v", err.Error()))
	}
	appConfig := container.GetAppConfig()

	errChan := make(chan error)
	go func() {
		errChan <- Run(root, appConfig.Port)
	}()

	go func() {
		errChan <- RunConsole()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errChan:
		color.Error.Printf("encounter error when starting server with [%s]\n", err.Error())
	case signal := <-signalChan:
		color.Warn.Printf("received signal [%s], process will exit\n", signal.String())
	}
	cancel()
	globalWg.Wait()
	container.Destroy()
}

// Run
//
//	@param ctx
func Run(ctx context.Context, port int64) error {
	router := gin.New()
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false

	for path, apiFunc := range container.GetNakedRouters() {
		router.Any(path, apiFunc)
	}

	router.Use(giner.DoResponseJSON())
	router.Use(container.GetMiddlewares()...)
	for path, apiFunc := range container.GetRouters() {
		router.Any(path, apiFunc)
	}
	router.Any("proxy/*api", gateway.Proxy)
	router.NoRoute(gateway.Handle)
	return router.Run(fmt.Sprintf(":%d", port))
}
