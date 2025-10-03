package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LSariol/Cove/internal/cli"
	"github.com/LSariol/Cove/internal/config"
	"github.com/LSariol/Cove/internal/server"
)

func main() {

	if err := config.Load(); err != nil {
		panic(err)
	}

	if err := config.Ensure(); err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go server.StartServer()

	cli.StartCLI()

	<-ctx.Done()
	log.Println("Shutting Down...")
}
