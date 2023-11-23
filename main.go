package main

import (
	"fmt"

	"github.com/crackeer/goaway/container"
	_ "github.com/crackeer/goaway/examples/sign"
	"github.com/crackeer/goaway/server"
	"github.com/gin-gonic/gin"
)

func init() {
	container.RegisterNakedRouter("access/token", func(ctx *gin.Context) {
		fmt.Println("access/token")
	})
}

func main() {
	server.Main()
}
