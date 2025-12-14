package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// AddIngredientHandler handles adding a new flavor or topping
func AddIngredientHandler(w http.ResponseWriter, r *http.Request) error {
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

	if r.Method == http.MethodPost {
		ingType := r.FormValue("ingredient_type")
		name := r.FormValue("name")
		img := r.FormValue("img")
		if ingType == "Flavor" {
			_ = models.AddFlavor(name, img)
		} else if ingType == "Topping" {
			price, _ := strconv.ParseFloat(r.FormValue("additional_price"), 64)
			_ = models.AddTopping(name, price, img)
		}
		http.Redirect(w, r, "/admin/ingredients", http.StatusSeeOther)
		return nil
	}

	tmpl := template.Must(template.ParseFiles(
		"views/ingredients/add-ingredient.gohtml",
		"views/layouts/layout.gohtml",
		"views/partials/head.gohtml",
		"views/partials/header.gohtml",
		"views/partials/navigation.gohtml",
		"views/partials/footer.gohtml",
	))
	data := struct {
		Title string
		User  *models.User
	}{
		Title: "Add Ingredient",
		User:  user,
	}
	return tmpl.ExecuteTemplate(w, "layout", data)
}
