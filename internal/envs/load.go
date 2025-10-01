package envs

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Load() error {

	// Load .env in dev
	if err := godotenv.Load(".env"); err == nil {
		return nil
	}

	if err := godotenv.Load("/app/vault/.env"); err == nil {
		return nil
	}

	return fmt.Errorf("no .env file found")
}
