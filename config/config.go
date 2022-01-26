package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

func Config(key string) string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	// load .env file
	err = godotenv.Load(filepath.Join(wd, ".env"))
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}
