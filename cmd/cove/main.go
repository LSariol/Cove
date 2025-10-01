package main

import (
	"github.com/LSariol/Cove/internal/cli"
	"github.com/LSariol/Cove/internal/envs"
	"github.com/LSariol/Cove/internal/server"
)

func main() {

	err := envs.Load()
	if err != nil {
		panic(err)
	}

	go server.StartServer()

	cli.StartCLI()
	// if os.Getenv("HEADLESS") != "true" {
	// 	cli.StartCLI()
	// } else {
	// 	select {}
	// }

}
