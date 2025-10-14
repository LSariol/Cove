package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LSariol/Cove/internal/database"
	"github.com/joho/godotenv"
)

const (
	USER = "cove_app"
	PSWD = "QJcwBUd4Buboau21wjki"
)

func main() {
	//Database testing here

	if err := godotenv.Load(".env"); err != nil {
		log.Panic("unable to load .env file")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db := database.NewDB()

	if err := db.Connect(ctx); err != nil {
		log.Panic("unable to connect to database")
	}

	// one := database.Secret{
	// 	Key:   "encrypted",
	// 	Value: "Secret",
	// }
	// err := db.CreateSecret(ctx, one)
	// if err != nil {
	// 	log.Panic(err)
	// }

	secrets, err := db.GetAllKeys(ctx)
	if err != nil {
		log.Println("get secret error: %w", err)
	}

	// s, err := db.GetSecret(ctx, "encrypted")
	// if err != nil {
	// 	log.Panic(err)
	// }
	// fmt.Println(s)

	// s.Value = "applesAndBannanas12345"
	// err = db.UpdateSecret(ctx, s)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// s, err = db.GetSecret(ctx, "foo")
	// if err != nil {
	// 	log.Panic(err)
	// }
	// fmt.Println(s)

	// err = db.DeleteSecret(ctx, "foo")
	// if err != nil {
	// 	log.Panic(err)
	// }

	// s, err = db.GetSecret(ctx, "foo")
	// if err != nil {
	// 	log.Panic(err)
	// }

	fmt.Println(secrets)

	<-ctx.Done()
	log.Println("Shutting Down...")

	db.Close()
	log.Println("database closed")

	stop()
}
