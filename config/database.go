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
	// Get environment variables or use defaults
	dbUser := getEnv("DB_USER", "anniedatabase_user")
	dbPassword := getEnv("DB_PASSWORD", "zS4wHtUYiDvwuhdAxPXgCmYngVT8SuED")
	dbHost := getEnv("DB_HOST", "dpg-d4fshta4d50c73f15jeg-a.oregon-postgres.render.com")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "anniedatabase")

	// Create the Postgres connection string
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
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
