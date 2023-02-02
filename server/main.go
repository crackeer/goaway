package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/crackeer/go-gateway/container"
	"github.com/gookit/color"
)

// Run a server
func Run() {
	container.InitConfig()
	container.InitGUID()

	run()
}

func run() {
	router := initEndpoint()
	errChan := make(chan error)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		errChan <- router.Run(fmt.Sprintf(":%d", container.Config.Port))
	}()

	/* ============================= Block here ============================= */
	select {
	case err := <-errChan:
		color.Error.Printf("encounter error when starting server with [%s]\n", err.Error())
	case signal := <-signalChan:
		color.Warn.Printf("received signal [%s], process will exit\n", signal.String())
	}

}
