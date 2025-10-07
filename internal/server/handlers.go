package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const boostrap bool = false

func defineRoutes(mux *http.ServeMux) {

	// Home route: no middleware needed
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/bootstrap/lighthouse", bootstrapHandler)

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

func bootstrapHandler(w http.ResponseWriter, r *http.Request) {

	// Create the marker file
	err := CreateBootstrapMarker()
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// Load client secret
	clientSecret := os.Getenv("COVE_CLIENT_SECRET")
	if clientSecret == "" {
		_ = DeleteBootstrapMarker()
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// Prepare JSON response
	response := struct {
		Secret string `json:"secret"`
	}{
		Secret: clientSecret,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		_ = DeleteBootstrapMarker()
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Bootstrap complete")
}

func CreateBootstrapMarker() error {
	dir := os.Getenv("APP_MARKER_PATH")
	if dir == "" {
		dir = "/app/vault/markers"
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create marker directory: %w", err)
	}

	marker := filepath.Join(dir, "bootstrap_completed")

	// Create marker file, fail if it already exists
	f, err := os.OpenFile(marker, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("failed to create marker file: %w", err)
	}
	_ = f.Close()

	return nil
}

// deleteBootstrapMarker removes the bootstrap marker file if it exists.
func DeleteBootstrapMarker() error {
	dir := os.Getenv("APP_MARKER_PATH")
	if dir == "" {
		dir = "/app/vault/markers"
	}

	marker := filepath.Join(dir, "bootstrap_completed")

	// Attempt to remove the file
	if err := os.Remove(marker); err != nil {
		// Ignore if file doesn't exist
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove marker file: %w", err)
		}
	}

	fmt.Println("bootstrap cleared")

	return nil
}
