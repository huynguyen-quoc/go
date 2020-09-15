package main

import (
	"github.com/huynguyen-quoc/go/echo/server"
)

func main() {
	serverRun := &server.Server{}
	serverRun.RunServer()
}
