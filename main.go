package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/crackeer/go-gateway/container"
	"github.com/crackeer/go-gateway/server"
	"github.com/gookit/color"
)

func main() {
	root := context.Background()
	gloabalWg := &sync.WaitGroup{}
	cancelCtx, cancel := context.WithCancel(root)
	if err := container.InitContainer(cancelCtx, gloabalWg); err != nil {
		panic(fmt.Sprintf("Failed to initialize container: %v", err.Error()))
	}
	appConfig := container.GetAppConfig()

	errChan := make(chan error)

	go func() {
		errChan <- server.Run(root, appConfig.Port)
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
	gloabalWg.Wait()
	container.Destroy()
}
