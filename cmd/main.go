package main

import (
	"github.com/g4ze/byoc/pkg/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env.postgres")
	if err != nil {
		panic(err)
	}
	routes.Server()
}
