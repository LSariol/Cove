package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LSariol/Cove/internal/cli"
	"github.com/LSariol/Cove/internal/envs"
	"github.com/LSariol/Cove/internal/server"
)

func main() {

	err := envs.Load()
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go server.StartServer()

	cli.StartCLI()

	<-ctx.Done()
	log.Println("Shutting Down...")
}
