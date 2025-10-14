package server

import "time"

// Return type of CreateSecret
type CreateResponse struct {
	Key       string    `json:"key"`
	DateAdded time.Time `json:"date_added"`
}

// Return type of GetPublicKeyVault
type PublicSecret struct {
	Key          string    `json:"key"`
	Version      int       `json:"version"`
	TimesPulled  int       `json:"timespulled"`
	DateAdded    time.Time `json:"dateAdded"`
	LastModified time.Time `json:"lastModified"`
}

//
type Secret struct {
	SecretID    string `json:"secretID"`
	SecretValue string `json:"secretValue"`
}
