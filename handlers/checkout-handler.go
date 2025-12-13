package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/ealbanca/Ice-Cream-Truck/config"
	"github.com/ealbanca/Ice-Cream-Truck/models"
)

func CheckoutHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		var orderID int
		var totalPrice float64
		var user *models.User
		var message, messageType string
		if userID, err := GetSessionUserID(r); err == nil {
			user, _ = models.GetUserByID(userID)
			if user != nil {
				err := config.DB.QueryRow("SELECT id FROM orders WHERE user_id = $1 AND status = 'cart'", user.ID).Scan(&orderID)
				if err == sql.ErrNoRows {
					err2 := config.DB.QueryRow("INSERT INTO orders (user_id, total_price, status) VALUES ($1, 0, 'cart') RETURNING id", user.ID).Scan(&orderID)
					if err2 != nil {
						message = "Failed to create order. Please try again."
						messageType = "error"
					}
				} else if err != nil {
					message = "Failed to find order. Please try again."
					messageType = "error"
				}
				totalPrice = models.GetCartTotal(user.ID)
			}
		} else {
			c, err := r.Cookie("guest_cart_id")
			if err == nil {
				guestCartID := c.Value
				err := config.DB.QueryRow("SELECT id FROM orders WHERE guest_cart_id = $1 AND status = 'cart'", guestCartID).Scan(&orderID)
				if err == sql.ErrNoRows {
					err2 := config.DB.QueryRow("INSERT INTO orders (guest_cart_id, total_price, status) VALUES ($1, 0, 'cart') RETURNING id", guestCartID).Scan(&orderID)
					if err2 != nil {
						message = "Failed to create order. Please try again."
						messageType = "error"
					}
				} else if err != nil {
					message = "Failed to find order. Please try again."
					messageType = "error"
				}
				guestCartItems := models.GetGuestCartItems(guestCartID)
				for _, item := range guestCartItems {
					totalPrice += item.TotalPrice
				}
			}
		}
		info := models.CheckoutInfo{
			OrderID:    orderID,
			FullName:   r.FormValue("full_name"),
			Email:      r.FormValue("email"),
			Phone:      r.FormValue("phone"),
			CreditCard: r.FormValue("credit_card"),
			TotalPrice: totalPrice,
		}
		err := models.SaveCheckoutInfo(info)
		if err != nil {
			message = "Failed to complete your order. Please try again."
			messageType = "error"
		} else {
			message = "Thank you for your order! You will be redirected to the home page."
			messageType = "success"
		}
		// Render thank you or error message and redirect if success
		tmpl := template.Must(template.ParseFiles(
			filepath.Join("views", "layouts", "layout.gohtml"),
			filepath.Join("views", "partials", "head.gohtml"),
			filepath.Join("views", "partials", "header.gohtml"),
			filepath.Join("views", "partials", "navigation.gohtml"),
			filepath.Join("views", "partials", "footer.gohtml"),
			filepath.Join("views", "cart", "checkout.gohtml"),
		))
		data := struct {
			Title       string
			Year        int
			OrderID     int
			User        *models.User
			TotalPrice  float64
			FullName    string
			Email       string
			Phone       string
			Message     string
			MessageType string
		}{
			Title:       "Checkout",
			Year:        2025,
			OrderID:     orderID,
			User:        user,
			TotalPrice:  totalPrice,
			FullName:    r.FormValue("full_name"),
			Email:       r.FormValue("email"),
			Phone:       r.FormValue("phone"),
			Message:     message,
			MessageType: messageType,
		}
		// If success, add JS redirect to home after 3 seconds
		if messageType == "success" {
			w.Header().Set("Refresh", "3; url=/")
		}
		err2 := tmpl.ExecuteTemplate(w, "layout", data)
		return err2
	}

	// For GET or on error, render the checkout page with OrderID
	orderID := 0 // TODO: set this to the current user's or guest's order ID
	if oid := r.URL.Query().Get("order_id"); oid != "" {
		orderID, _ = strconv.Atoi(oid)
	}
	var user *models.User
	var totalPrice float64
	if userID, err := GetSessionUserID(r); err == nil {
		user, _ = models.GetUserByID(userID)
		if user != nil {
			totalPrice = models.GetCartTotal(user.ID)
		}
	} else {
		// Guest: check for guest_cart_id cookie
		c, err := r.Cookie("guest_cart_id")
		if err == nil {
			guestCartID := c.Value
			guestCartItems := models.GetGuestCartItems(guestCartID)
			for _, item := range guestCartItems {
				totalPrice += item.TotalPrice
			}
		}
	}

	fullName := ""
	email := ""
	phone := ""
	if user != nil {
		fullName = user.Name + " " + user.Lastname
		email = user.Email
		phone = user.Phone
	}

	tmpl := template.Must(template.ParseFiles(
		filepath.Join("views", "layouts", "layout.gohtml"),
		filepath.Join("views", "partials", "head.gohtml"),
		filepath.Join("views", "partials", "header.gohtml"),
		filepath.Join("views", "partials", "navigation.gohtml"),
		filepath.Join("views", "partials", "footer.gohtml"),
		filepath.Join("views", "cart", "checkout.gohtml"),
	))
	data := struct {
		Title       string
		Year        int
		OrderID     int
		User        *models.User
		TotalPrice  float64
		FullName    string
		Email       string
		Phone       string
		Message     string
		MessageType string
	}{
		Title:      "Checkout",
		Year:       2025, // or use time.Now().Year()
		OrderID:    orderID,
		User:       user,
		TotalPrice: totalPrice,
		FullName:   fullName,
		Email:      email,
		Phone:      phone,
	}
	err2 := tmpl.ExecuteTemplate(w, "layout", data)
	return err2
}
