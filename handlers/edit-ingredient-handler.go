package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// EditIngredientHandler displays the edit form for a flavor or topping
func EditIngredientHandler(w http.ResponseWriter, r *http.Request) error {
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
	idStr := r.URL.Path[len("/ingredients/edit/"):]
	ingredientID, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return nil
	}
	ingredient, _ := models.GetIngredientByID(ingredientID)
	tmpl := template.Must(template.ParseFiles(
		"views/ingredients/edit-ingredient.gohtml",
		"views/layouts/layout.gohtml",
		"views/partials/head.gohtml",
		"views/partials/header.gohtml",
		"views/partials/navigation.gohtml",
		"views/partials/footer.gohtml",
	))
	data := struct {
		Title      string
		User       *models.User
		Ingredient models.Ingredient
	}{
		Title:      "Edit Ingredient",
		User:       user,
		Ingredient: ingredient,
	}
	return tmpl.ExecuteTemplate(w, "layout", data)
}
