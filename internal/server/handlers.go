package server

import (
	"fmt"
	"net/http"
)

func defineRoutes() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/token", tokenHandler)
	http.HandleFunc("/store", storeSecretHandler)
	http.HandleFunc("/get", getSecretHandler)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "token Handler")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Keyvault")
}

func storeSecretHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "store handler")
}

func getSecretHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "get handler")
}
