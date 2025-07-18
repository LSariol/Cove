package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func defineRoutes(mux *http.ServeMux) {

	// Home route: no middleware needed
	mux.HandleFunc("/health", healthHandler)

	// Routes with authentication middleware
	mux.Handle("/secrets", authenticateClientSecret(http.HandlerFunc(secretHandler)))
	mux.Handle("/secrets/", authenticateClientSecret(http.HandlerFunc(secretHandler)))
	mux.Handle("/auth", authenticateClientSecret(http.HandlerFunc(authHandler)))

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Healthy bool   `json:"healthy"`
		Time    string `json:"time"`
	}{
		Healthy: true,
		Time:    time.Now().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Health Check Successful")

}

func authHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Authenticated bool   `json:"Authenticated"`
		Time          string `json:"time"`
	}{
		Authenticated: true,
		Time:          time.Now().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Authentication Check Successful")

}

func secretHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/secrets" {
		if r.Method == http.MethodGet {
			getAllSecrets(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/secrets/")
	if id == "" {
		http.Error(w, "Missing secret ID", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:

		getSecret(w, r, id)

	case http.MethodPost:

		postSecret(w, r, id)

	case http.MethodDelete:

		deleteSecret(w, r, id)

	case http.MethodPatch:

		patchSecret(w, r, id)

	default:

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}
