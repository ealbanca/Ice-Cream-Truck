package handlers

import (
	"net/http"
	"strconv"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// DeleteIngredientHandler deletes a flavor or topping by ID
func DeleteIngredientHandler(w http.ResponseWriter, r *http.Request) error {
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
	idStr := r.URL.Path[len("/ingredients/delete/"):]
	ingredientID, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return nil
	}
	_ = models.DeleteIngredientByID(ingredientID)
	http.Redirect(w, r, "/admin/ingredients", http.StatusSeeOther)
	return nil
}
