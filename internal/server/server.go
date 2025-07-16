package server

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer() {

	//multiplexer (router)
	mux := http.NewServeMux()
	defineRoutes(mux)

	fmt.Println("Server is running at http://localhost:8081")
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", mux))
}
