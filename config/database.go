package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	// Get environment variables (all must be set)
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	missing := ""
	if dbUser == "" {
		missing += "DB_USER "
	}
	if dbPassword == "" {
		missing += "DB_PASSWORD "
	}
	if dbHost == "" {
		missing += "DB_HOST "
	}
	if dbPort == "" {
		missing += "DB_PORT "
	}
	if dbName == "" {
		missing += "DB_NAME "
	}
	if missing != "" {
		log.Fatalf("Missing required environment variables: %s", missing)
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open database connection
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
}

// Helper function to get environment variables
func getEnv(key, fallback string) string {
	// Deprecated: All credentials must come from environment variables only.
	return ""
}
