package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func GetLB_DNS(subdomain string) (string, error) {
	// Connection parameters
	err := godotenv.Load("../.env.postgres")
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}
	var (
		host     = "localhost"
		port     = 5432 // Default PostgreSQL port
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = os.Getenv("POSTGRES_DB")
	)

	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return "", fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		return "", fmt.Errorf("error pinging database: %v", err)
	}

	// Query the database
	var serviceURL string
	query := "SELECT \"loadbalancerDNS\" FROM \"Service\" WHERE \"Slug\" = $1"
	err = db.QueryRow(query, subdomain).Scan(&serviceURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no service URL found for subdomain: %s", subdomain)
		}
		return "", fmt.Errorf("error querying database: %v", err)
	}

	return serviceURL, nil
}
