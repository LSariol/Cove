package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Pulls client secret from environment variables
func getClientSecret() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("COVE_CLIENT_SECRET")
}

// Middleware function to handle authentication of client secret
func authenticateClientSecret(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing client secret", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Unauthorized: Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		providedSecret := tokenParts[1]
		storedSecret := getClientSecret()

		if providedSecret != storedSecret {
			http.Error(w, "Unauthorized: Invalid client secret", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)

	})
}
