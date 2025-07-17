package server

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer() {

	fmt.Println("Starting Server")

	//multiplexer (router)
	mux := http.NewServeMux()
	defineRoutes(mux)

	fmt.Println("Routes Defined")

	fmt.Println("Server is running at http://localhost:8081")
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", mux))
}
