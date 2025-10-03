package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func StartServer() {

	fmt.Println("Starting Server")

	//multiplexer (router)
	mux := http.NewServeMux()
	defineRoutes(mux)

	fmt.Println("Routes Defined")

	port := os.Getenv("APP_PORT")
	address := "0.0.0.0:" + port

	fmt.Println("Server is running at http://localhost:" + port)

	log.Fatal(http.ListenAndServe(address, mux))
}
