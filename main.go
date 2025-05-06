package main

import (
	"log"
	"os"
	"routinist/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	app.Run()
}
