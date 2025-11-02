package config

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "log"
    "os"
)

var DB *sql.DB

func InitDB() {
    // Get environment variables or use defaults
    dbUser := getEnv("DB_USER", "root")
    dbPassword := getEnv("DB_PASSWORD", "")
    dbHost := getEnv("DB_HOST", "localhost")
    dbPort := getEnv("DB_PORT", "3306")
    dbName := getEnv("DB_NAME", "icecream_truck")

    // Create the connection string
    connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        dbUser, dbPassword, dbHost, dbPort, dbName)

    // Open database connection
    var err error
    DB, err = sql.Open("mysql", connectionString)
    if err != nil {
        log.Fatal("Error connecting to the database: ", err)
    }

    // Test the connection
    err = DB.Ping()
    if err != nil {
        log.Fatal("Error pinging the database: ", err)
    }

    log.Println("Successfully connected to MySQL database")
}

// Helper function to get environment variables
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}