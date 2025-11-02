package main

import (
	"fmt"

	"github.com/ealbanca/Ice-Cream-Truck/config"
	"github.com/ealbanca/Ice-Cream-Truck/models"
)

func main() {
	// Initialize database connection
	config.InitDB()
	defer config.DB.Close()

	// Test connection
	err := config.DB.Ping()
	if err != nil {
		fmt.Printf("❌ Database connection failed: %v\n", err)
		return
	}
	fmt.Println("✅ Successfully connected to database")

	// Try to get all products
	products, err := models.GetAllProducts()
	if err != nil {
		fmt.Printf("❌ Failed to query products: %v\n", err)
		return
	}
	fmt.Println("✅ Successfully queried products table")
	fmt.Printf("Found %d products in database\n", len(products))

	// Test creating a product
	testProduct := &models.Product{
		Name:        "Test Ice Cream",
		Description: "Test product - please ignore",
		Price:       1.99,
	}

	err = testProduct.Create()
	if err != nil {
		fmt.Printf("❌ Failed to create test product: %v\n", err)
		return
	}
	fmt.Printf("✅ Successfully created test product with ID: %d\n", testProduct.ID)
}
