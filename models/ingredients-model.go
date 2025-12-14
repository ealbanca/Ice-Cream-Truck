package models

import (
	"database/sql"

	"github.com/ealbanca/Ice-Cream-Truck/config"
)

type Ingredient struct {
	ID    int
	Name  string
	Type  string  // "Flavor" or "Topping"
	Img   string  // For flavor_img or topping_img
	Price float64 // For additional_price (topping only)
}

// GetAllIngredients returns all flavors and toppings as ingredients
func GetAllIngredients() ([]Ingredient, error) {
	var ingredients []Ingredient

	// Flavors
	rows, err := config.DB.Query("SELECT id, flavor_name FROM flavors")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var i Ingredient
			if err := rows.Scan(&i.ID, &i.Name); err == nil {
				i.Type = "Flavor"
				ingredients = append(ingredients, i)
			}
		}
	}

	// Toppings
	rows, err = config.DB.Query("SELECT id, topping_name FROM toppings")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var i Ingredient
			if err := rows.Scan(&i.ID, &i.Name); err == nil {
				i.Type = "Topping"
				ingredients = append(ingredients, i)
			}
		}
	}

	return ingredients, nil
}

// GetIngredientByID fetches a flavor or topping by ID and returns as Ingredient
func GetIngredientByID(id int) (Ingredient, error) {
	var i Ingredient
	// Try flavor first
	row := config.DB.QueryRow("SELECT id, flavor_name, falvor_img FROM flavors WHERE id = $1", id)
	var img sql.NullString
	err := row.Scan(&i.ID, &i.Name, &img)
	if err == nil {
		i.Type = "Flavor"
		i.Img = img.String
		return i, nil
	}
	// Try topping
	row = config.DB.QueryRow("SELECT id, topping_name, additional_price, topping_img FROM toppings WHERE id = $1", id)
	var price sql.NullFloat64
	err = row.Scan(&i.ID, &i.Name, &price, &img)
	if err == nil {
		i.Type = "Topping"
		i.Price = price.Float64
		i.Img = img.String
		return i, nil
	}
	return i, err
}
