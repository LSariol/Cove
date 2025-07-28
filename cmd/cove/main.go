package main

import (
	"os"

	"github.com/LSariol/Cove/internal/cli"
	"github.com/LSariol/Cove/internal/server"
)

func main() {

	server.StartServer()

	if os.Getenv("HEADLESS") != "true" {
		cli.StartCLI()
	} else {
		select {}
	}

}
