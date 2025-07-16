package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LSariol/Cove/internal/encryption"
)

type payload struct {
	SecretID    string `json:"secretID"`
	SecretValue string `json:"secretValue"`
}

type response struct {
	Message string `json:"message"`
}

func getAllSecrets(w http.ResponseWriter, r *http.Request) {
	publicKeyVault := encryption.GetPublicVault()
	fmt.Println(publicKeyVault)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(publicKeyVault); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Cove - getAllSecrets: failed to encode response:", err)
		return
	}

	fmt.Println("GetAllSecrets Complete")

}

func getSecret(w http.ResponseWriter, r *http.Request, id string) {

	secret, err := encryption.GetSecret(id)
	if !err {
		http.Error(w, secret, http.StatusInternalServerError)
		fmt.Println(secret)
	}

	response := struct {
		Secret string `json:"secret"`
	}{
		Secret: secret,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Cove - getAllSecrets: failed to encode response:", err)
		return
	}

	fmt.Println("GetSecret Complete")

}

func postSecret(w http.ResponseWriter, r *http.Request, ID string) {

	fmt.Println("Cove - postSecret")

	var load payload

	if err := json.NewDecoder(r.Body).Decode(&load); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		fmt.Println("Cove - postSecret: failed to decode JSON:", err)
		return
	}

	msg, ok := encryption.AddSecret(ID, load.SecretValue)
	if !ok {
		http.Error(w, msg, http.StatusBadRequest)
		fmt.Println("Cove - postSecret: ", msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := response{
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

	fmt.Println("PostSecret Complete")

}

func patchSecret(w http.ResponseWriter, r *http.Request, ID string) {

	fmt.Println("Cove - patchSecret")
	var load payload

	if err := json.NewDecoder(r.Body).Decode(&load); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		fmt.Println("Cove - patchSecret: failed to decode JSON:", err)
		return
	}

	msg, ok := encryption.UpdateSecret(ID, load.SecretValue)
	if !ok {
		http.Error(w, msg, http.StatusBadRequest)
		fmt.Println("Cove - patchSecret: ", msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := response{
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

	fmt.Println("Patch Secret Complete")

}

func deleteSecret(w http.ResponseWriter, r *http.Request, ID string) {
	fmt.Println("Cove - deleteSecret")

	msg, ok := encryption.RemoveSecret(ID)
	if !ok {
		http.Error(w, msg, http.StatusBadRequest)
		fmt.Println("Cove - deleteSecret: ", msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := response{
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

	fmt.Println("DeleteSecret Complete")

}
