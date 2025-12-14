package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// AdminOrdersHandler serves the admin all orders page
func AdminOrdersHandler(w http.ResponseWriter, r *http.Request) error {
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
	orders, err := models.GetAllOrders()
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return nil
	}
	tmpl, err := template.ParseFiles(
		"views/order/orders.gohtml",
		"views/layouts/layout.gohtml",
		"views/partials/head.gohtml",
		"views/partials/header.gohtml",
		"views/partials/navigation.gohtml",
		"views/partials/footer.gohtml",
	)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return nil
	}
	data := struct {
		Title  string
		Year   int
		User   *models.User
		Orders []models.Order
	}{
		Title:  "All Orders",
		Year:   time.Now().Year(),
		User:   user,
		Orders: orders,
	}
	tmpl.ExecuteTemplate(w, "layout", data)
	return nil
}

func OrderDetailsHandler(w http.ResponseWriter, r *http.Request) error {
	// Parse order ID from URL: /order/{id}
	idStr := r.URL.Path[len("/order/"):] // crude parse, assumes /order/{id}
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return nil
	}

	userID, err := GetSessionUserID(r)
	if err != nil || userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil
	}

	   order, err := models.GetOrderByID(orderID)
	   if err != nil || (!order.UserID.Valid || order.UserID.Int64 != int64(userID)) {
		   http.NotFound(w, r)
		   return nil
	   }
	items, err := models.GetOrderItems(orderID)
	if err != nil {
		http.Error(w, "Failed to fetch order items", http.StatusInternalServerError)
		return nil
	}
	user, _ := models.GetUserByID(userID)

	tmpl, err := template.ParseFiles(
		"views/order/order-details.gohtml",
		"views/layouts/layout.gohtml",
		"views/partials/head.gohtml",
		"views/partials/header.gohtml",
		"views/partials/navigation.gohtml",
		"views/partials/footer.gohtml",
	)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return nil
	}

	data := struct {
		Title string
		Year  int
		User  *models.User
		Order *models.Order
		Items []models.OrderItem
	}{
		Title: "Order Details",
		Year:  time.Now().Year(),
		User:  user,
		Order: order,
		Items: items,
	}
	tmpl.ExecuteTemplate(w, "layout", data)
	return nil
}

// OrdersHandler serves the user's orders page
func OrdersHandler(w http.ResponseWriter, r *http.Request) error {
	userID, err := GetSessionUserID(r)
	if err != nil || userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil
	}

	user, _ := models.GetUserByID(userID)
	orders, err := models.GetOrdersByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return nil
	}

	tmpl, err := template.ParseFiles(
		"views/order/orders.gohtml",
		"views/layouts/layout.gohtml",
		"views/partials/head.gohtml",
		"views/partials/header.gohtml",
		"views/partials/navigation.gohtml",
		"views/partials/footer.gohtml",
	)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return nil
	}

	data := struct {
		Title  string
		Year   int
		User   *models.User
		Orders []models.Order
	}{
		Title:  "Your Orders",
		Year:   time.Now().Year(),
		User:   user,
		Orders: orders,
	}
	tmpl.ExecuteTemplate(w, "layout", data)
	return nil
}
