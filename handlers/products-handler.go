package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		products, err := models.GetAllProducts()
		if err != nil {
			return fmt.Errorf("failed to get products: %w", err)
		}
		json.NewEncoder(w).Encode(products)
		return nil

	case http.MethodPost:
		var product models.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			return fmt.Errorf("invalid product input: %w", err)
		}

		if err := product.Create(); err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
		return nil

	default:
		return fmt.Errorf("method not allowed")
	}
}
