package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LSariol/Cove/internal/database"
)

type Server struct {
	DB *database.Database
}

func NewServer(db *database.Database) *Server {
	return &Server{
		DB: db,
	}
}

func (s *Server) Start() {

	//multiplexer (router)
	mux := http.NewServeMux()
	s.defineRoutes(mux)

	port := os.Getenv("APP_PORT")
	address := "0.0.0.0:" + port

	fmt.Printf("Running on %s\n", address)

	log.Fatal(http.ListenAndServe(address, mux))
}
