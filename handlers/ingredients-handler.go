package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// ManageIngredientsHandler displays all ingredients (flavors and toppings)
func ManageIngredientsHandler(w http.ResponseWriter, r *http.Request) error {
	userID, err := GetSessionUserID(r)
	if err != nil || userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil
	}
	user, _ := models.GetUserByID(userID)
	if user == nil || user.UserType != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return nil
	}
	ingredients, _ := models.GetAllIngredients()
	tmpl := template.Must(template.ParseFiles(
		"views/ingredients/manage-ingredients.gohtml",
		"views/layouts/layout.gohtml",
		"views/partials/head.gohtml",
		"views/partials/header.gohtml",
		"views/partials/navigation.gohtml",
		"views/partials/footer.gohtml",
	))
	data := struct {
		Title       string
		Year        int
		User        *models.User
		Ingredients []models.Ingredient
	}{
		Title:       "Manage Ingredients",
		Year:        time.Now().Year(),
		User:        user,
		Ingredients: ingredients,
	}
	return tmpl.ExecuteTemplate(w, "layout", data)
}
