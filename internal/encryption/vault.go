package encryption

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type secretEntry struct {
	Key          string `json:"key"`
	Secret       string `json:"secret"`
	DateAdded    string `json:"dateAdded"`
	LastModified string `json:"lastModified"`
	Version      string `json:"version"`
}

type PublicSecretEntry struct {
	Key          string `json:"key"`
	DateAdded    string `json:"dateAdded"`
	LastModified string `json:"lastModified"`
	Version      string `json:"version"`
}

func GetPublicVault() []PublicSecretEntry {
	keyVault := readJSON()
	var publicVault []PublicSecretEntry

	for _, secret := range keyVault {
		publicEntry := PublicSecretEntry{
			Key:          secret.Key,
			DateAdded:    secret.DateAdded,
			LastModified: secret.LastModified,
			Version:      secret.Version,
		}

		publicVault = append(publicVault, publicEntry)
	}

	return publicVault

}

func GetSecret(key string) (string, bool) {
	keyVault := readJSON()

	for _, entry := range keyVault {
		if entry.Key == key {
			secret, err := decryptData(entry.Secret)
			if err != nil {
				log.Fatal(err)
			}
			return secret, true
		}
	}

	return "Secret does not exist", false

}

func AddSecret(name string, value string) (string, bool) {
	keyVault := readJSON()

	if exists(name, keyVault) {
		return "A secret with that name already exists.", false
	}

	encryptedValue, err := encryptData(value)
	if err != nil {
		log.Fatal(err)
	}

	var entry secretEntry
	entry.Key = name
	entry.Secret = encryptedValue
	entry.DateAdded = time.Now().Format("2006-01-02 15:04:05")
	entry.LastModified = time.Now().Format("2006-01-02 15:04:05")
	entry.Version = "1"

	keyVault = append(keyVault, entry)

	storeJSON(keyVault)

	return "Secret has been created", true

}

func RemoveSecret(name string) (string, bool) {
	keyVault := readJSON()

	indexToRemove := -1

	for index, entry := range keyVault {
		if entry.Key == name {
			indexToRemove = index
			break
		}
	}

	if indexToRemove == -1 {
		return "Cove - removeSecret: Secret doesnt exist. Unable to remove", false
	}

	keyVault = append(keyVault[:indexToRemove], keyVault[indexToRemove+1:]...)
	storeJSON(keyVault)

	return "Cove - removeSecret: " + name + " has been removed.", true

}

func UpdateSecret(name string, newValue string) (string, bool) {
	keyVault := readJSON()

	updated := false

	for i := range keyVault {
		if keyVault[i].Key == name {
			secret, err := encryptData(newValue)
			if err != nil {
				return "Cove - UpdateSecret: Error during encryption", false
			}
			keyVault[i].Secret = secret
			keyVault[i].LastModified = time.Now().Format("2006-01-02 15:04:05")
			keyVault[i].Version = incrementVersion(keyVault[i].Version)
			updated = true
			break
		}
	}

	if !updated {
		return "Cove - UpdateSecret: Secret does not exist. It cannot be updated.", false
	}

	storeJSON(keyVault)
	return "Cove - UpdateSecret: Secret has been updated.", true

}

func exists(name string, keyVault []secretEntry) bool {

	for _, entry := range keyVault {
		if entry.Key == name {
			return true
		}
	}

	return false

}

func readJSON() []secretEntry {

	var keyVault []secretEntry

	data, err := os.ReadFile("../../internal/encryption/vault.json")
	if err != nil {
		fmt.Println("Cove - Vault: Failed to load vault.json")
		fmt.Println(err)
	}

	//Unmarshal repos.json into watchList
	err = json.Unmarshal([]byte(data), &keyVault)
	if err != nil {
		fmt.Println("Cove - Vault: Failed to Unmarshal json into WatchedRepos.")
		fmt.Println(err)
	}

	return keyVault
}

func storeJSON(keyVault []secretEntry) {

	keyVaultData, err := json.MarshalIndent(keyVault, "", "	")
	if err != nil {
		fmt.Println("Cove - storeJSON: Failed to Marhsal json into keyVault.")
	}

	err = os.WriteFile("../../internal/encryption/vault.json", keyVaultData, 0644)
	if err != nil {
		fmt.Println("Cove - storeJSON: Failed to write to vault.json.")
	}

	fmt.Println("Cove - StoreJSON: Keyvault has been stored.")
}

func incrementVersion(version string) string {
	v, _ := strconv.Atoi(version)
	v++
	return strconv.Itoa(v)
}
