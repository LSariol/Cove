package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func getClientSecret() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func StartServer() {

	defineRoutes()
	fmt.Println("Server is running at http://localhost:8081")
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", nil))
}
