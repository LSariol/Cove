package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LSariol/Cove/internal/database"
)

// Get benign vault data for CLI
func (s *Server) getAllSecrets(w http.ResponseWriter, r *http.Request) {

	keys, err := s.DB.GetAllKeys(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var pubKeys []PublicSecret
	for _, key := range keys {
		var pubKey PublicSecret
		pubKey.Key = key.Key
		pubKey.Version = key.Version
		pubKey.TimesPulled = key.TimesPulled
		pubKey.DateAdded = key.DateAdded
		pubKey.LastModified = key.LastModified
		pubKeys = append(pubKeys, pubKey)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(pubKeys); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Cove - getAllSecrets: failed to encode response:", err)
		return
	}

	log.Println("All secrets got.")

}

func (s *Server) getSecret(w http.ResponseWriter, r *http.Request, ID string) {

	secret, err := s.DB.GetSecret(r.Context(), ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response := packPayload(secret)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to get", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	log.Printf("%s got.\n", ID)

}

func (s *Server) postSecret(w http.ResponseWriter, r *http.Request, ID string) {

	fmt.Println("Cove - postSecret")

	var secret Secret

	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		fmt.Println("Cove - postSecret: failed to decode JSON:", err)
		return
	}

	var dbSecret database.Secret = database.Secret{
		Key:   ID,
		Value: secret.SecretValue,
	}

	secre, err := s.DB.CreateSecret(r.Context(), dbSecret)
	if err != nil {
		http.Error(w, "failed to create", http.StatusBadRequest)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	res := packPayload(secre)

	json.NewEncoder(w).Encode(res)

	log.Printf("%s has been created.\n", ID)

}

func (s *Server) patchSecret(w http.ResponseWriter, r *http.Request, ID string) {

	fmt.Println("Cove - patchSecret")
	var secret Secret

	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		fmt.Println("Cove - patchSecret: failed to decode JSON:", err)
		return
	}

	var dbSecret database.Secret = database.Secret{
		Key:   ID,
		Value: secret.SecretValue,
	}

	err := s.DB.UpdateSecret(r.Context(), dbSecret)
	if err != nil {
		http.Error(w, "patch failed", http.StatusBadRequest)
		log.Println(err)
		return
	}

	// No body 204 response
	w.WriteHeader(http.StatusNoContent)

	log.Printf("%s has been patched.\n", ID)

}

func (s *Server) deleteSecret(w http.ResponseWriter, r *http.Request, ID string) {
	fmt.Println("Cove - deleteSecret")

	err := s.DB.DeleteSecret(r.Context(), ID)
	if err != nil {
		http.Error(w, "delete failed", http.StatusBadRequest)
		log.Println(err)
		return
	}

	// No body 204 repsonse.
	w.WriteHeader(http.StatusNoContent)
	log.Printf("%s has been deleted.\n", ID)

}

func packPayload(s database.Secret) Secret {
	var secret Secret

	secret.SecretID = s.Key
	secret.SecretValue = s.Value

	return secret
}
