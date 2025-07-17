package main

import (
	"time"

	"github.com/LSariol/Cove/internal/cli"
	"github.com/LSariol/Cove/internal/server"
)

func main() {

	go server.StartServer()

	time.Sleep(2 * time.Second)
	cli.StartCLI()

}
