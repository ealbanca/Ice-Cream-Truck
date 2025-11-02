package main

import (
    "log"
    "net/http"
    "github.com/ealbanca/Ice-Cream-Truck/config"
    "github.com/ealbanca/Ice-Cream-Truck/handlers"
)

func main() {
    // Initialize database connection
    config.InitDB()
    defer config.DB.Close()

    // Initialize routes
    http.HandleFunc("/api/products", handlers.ProductsHandler)
    
    // Start the server
    log.Println("Server starting on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}