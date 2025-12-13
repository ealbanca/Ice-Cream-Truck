package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// UpdateCartQuantityHandler updates the quantity of a cart item
func UpdateCartQuantityHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		idStr := r.FormValue("cart_item_id")
		qtyStr := r.FormValue("quantity")
		id, _ := strconv.Atoi(idStr)
		qty, _ := strconv.Atoi(qtyStr)
		if id > 0 && qty > 0 {
			// Check if user is logged in
			if userID, err := GetSessionUserID(r); err == nil && userID > 0 {
				models.UpdateCartItemQuantity(id, qty)
			} else {
				// Guest: update guest cart item
				models.UpdateGuestCartItemQuantity(id, qty)
			}
		}
	}
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
	return nil
}

func CartHandler(w http.ResponseWriter, r *http.Request) error {
	var user *models.User
	var cartItems []models.CartItem
	var guestCartItems []models.GuestCartItem
	var cartTotal float64
	var guestCartID string
	if userID, err := GetSessionUserID(r); err == nil {
		user, _ = models.GetUserByID(userID)
		if user != nil {
			cartItems = models.GetCartItems(user.ID)
			cartTotal = models.GetCartTotal(user.ID)
		}
	} else {
		// Guest: check for guest_cart_id cookie
		c, err := r.Cookie("guest_cart_id")
		if err == nil {
			guestCartID = c.Value
			guestCartItems = models.GetGuestCartItems(guestCartID)
			for _, item := range guestCartItems {
				cartTotal += item.TotalPrice
			}
		}
	}
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
		Title          string
		Year           int
		User           *models.User
		CartItems      []models.CartItem
		GuestCartItems []models.GuestCartItem
		CartTotal      float64
		Message        string
		MessageType    string
	}{
		Title:          "Your Cart",
		Year:           time.Now().Year(),
		User:           user,
		CartItems:      cartItems,
		GuestCartItems: guestCartItems,
		CartTotal:      cartTotal,
		Message:        message,
		MessageType:    messageType,
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
		// Check if user is logged in
		if userID, err := GetSessionUserID(r); err == nil && userID > 0 {
			models.RemoveFromCart(id)
		} else {
			// Guest: remove from guest cart
			models.RemoveFromGuestCart(id)
		}
	}
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
	return nil
}
