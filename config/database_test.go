package config

import (
	"os"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	// Set test environment variables
	os.Setenv("DB_USER", "root")         // Replace with your MySQL username
	os.Setenv("DB_PASSWORD", "password") // Replace with your MySQL password
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "icecream_truck")

	// Initialize database connection
	InitDB()
	if DB == nil {
		t.Fatal("Database connection is nil")
	}

	// Test the connection with a ping
	err := DB.Ping()
	if err != nil {
		t.Fatalf("Could not ping database: %v", err)
	}

	// Try a simple query
	_, err = DB.Query("SELECT 1")
	if err != nil {
		t.Fatalf("Could not execute test query: %v", err)
	}

	t.Log("Successfully connected to the database and executed test query")
}
