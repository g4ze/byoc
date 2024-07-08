package main

import (
	"log"
	"os"

	"github.com/g4ze/byoc/pkg/routes"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
	routes.Server()
}
