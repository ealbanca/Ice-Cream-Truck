package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// Get all products
		products, err := models.GetAllProducts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(products)

	case http.MethodPost:
		// Create a new product
		var product models.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := product.Create(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
