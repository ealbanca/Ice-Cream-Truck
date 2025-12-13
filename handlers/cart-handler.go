package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

func CartHandler(w http.ResponseWriter, r *http.Request) error {
	var user *models.User
	if userID, err := GetSessionUserID(r); err == nil {
		user, _ = models.GetUserByID(userID)
	}
	userID := 0
	if user != nil {
		userID = user.ID
	}
	cartItems := models.GetCartItems(userID)
	cartTotal := models.GetCartTotal()
	message := ""
	messageType := ""
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("views", "layouts", "layout.gohtml"),
		filepath.Join("views", "partials", "head.gohtml"),
		filepath.Join("views", "partials", "header.gohtml"),
		filepath.Join("views", "partials", "navigation.gohtml"),
		filepath.Join("views", "partials", "footer.gohtml"),
		filepath.Join("views", "cart", "cart.gohtml"),
	))
	data := struct {
		Title       string
		Year        int
		User        *models.User
		CartItems   []models.CartItem
		CartTotal   float64
		Message     string
		MessageType string
	}{
		Title:       "Your Cart",
		Year:        time.Now().Year(),
		User:        user,
		CartItems:   cartItems,
		CartTotal:   cartTotal,
		Message:     message,
		MessageType: messageType,
	}
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		return err
	}
	return nil
}

func RemoveFromCartHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		idStr := r.FormValue("cart_item_id")
		id, _ := strconv.Atoi(idStr)
		models.RemoveFromCart(id)
	}
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
	return nil
}
