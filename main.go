package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/config"
	"github.com/ealbanca/Ice-Cream-Truck/handlers"
	"github.com/ealbanca/Ice-Cream-Truck/models"
	"github.com/joho/godotenv"
)

func main() {
	// Admin all orders page
	http.HandleFunc("/admin/orders", handlers.ErrorHandler(handlers.AdminOrdersHandler))
	// Load environment variables from .env file (for local development)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file (this is OK in production)")
	}

	// Initialize database connection
	config.InitDB()
	defer config.DB.Close()

	// Initialize routes
	http.HandleFunc("/api/products", handlers.ErrorHandler(handlers.ProductsHandler))
	http.HandleFunc("/api/events", handlers.ErrorHandler(handlers.EventHandler))
	http.HandleFunc("/about", handlers.ErrorHandler(handlers.AboutHandler))
	http.HandleFunc("/contact", handlers.ErrorHandler(handlers.ContactHandler))
	http.HandleFunc("/api/contact", handlers.ErrorHandler(handlers.ContactHandler))
	http.HandleFunc("/register", handlers.ErrorHandler(handlers.RegisterHandler))
	http.HandleFunc("/login", handlers.ErrorHandler(handlers.LoginHandler))
	http.HandleFunc("/logout", handlers.ErrorHandler(handlers.LogoutHandler))
	http.HandleFunc("/edit-account", handlers.ErrorHandler(handlers.EditAccountHandler))
	http.HandleFunc("/management", handlers.ErrorHandler(handlers.LoginHandler))
	http.HandleFunc("/build", handlers.ErrorHandler(handlers.BuildHandler))
	http.HandleFunc("/cart", handlers.ErrorHandler(handlers.CartHandler))
	http.HandleFunc("/cart/remove", handlers.ErrorHandler(handlers.RemoveFromCartHandler))
	http.HandleFunc("/cart/update-qty", handlers.ErrorHandler(handlers.UpdateCartQuantityHandler))
	http.HandleFunc("/checkout", handlers.ErrorHandler(handlers.CheckoutHandler))
	http.HandleFunc("/account/orders", handlers.ErrorHandler(handlers.OrdersHandler))
	http.HandleFunc("/account/creations", handlers.ErrorHandler(handlers.CreationsHandler))
	http.HandleFunc("/admin/ingredients", handlers.ErrorHandler(handlers.ManageIngredientsHandler))
	// Route for editing an ingredient
	http.HandleFunc("/ingredients/edit/", handlers.ErrorHandler(handlers.EditIngredientHandler))
	http.HandleFunc("/creation/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > len("/creation/") && r.Method == http.MethodGet && len(r.URL.Path) > len("/creation/") && r.URL.Path[len(r.URL.Path)-len("/delete"):] == "/delete" {
			handlers.ErrorHandler(handlers.CreationDeleteHandler)(w, r)
			return
		}
		http.NotFound(w, r)
	})
	// Order details page
	http.HandleFunc("/order/", handlers.ErrorHandler(handlers.OrderDetailsHandler))

	// Serve static files from public directory
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("public/js"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("public/images"))))

	// Parse templates (layout, partials, index)
	templates := template.Must(template.ParseFiles(
		filepath.Join("views", "layouts", "layout.gohtml"),
		filepath.Join("views", "partials", "head.gohtml"),
		filepath.Join("views", "partials", "header.gohtml"),
		filepath.Join("views", "partials", "navigation.gohtml"),
		filepath.Join("views", "partials", "footer.gohtml"),
		filepath.Join("views", "index.gohtml"),
	))

	// Root route handler (renders homepage)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		if userID, err := handlers.GetSessionUserID(r); err == nil {
			user, _ = models.GetUserByID(userID)
		}
		data := struct {
			Title string
			Year  int
			User  *models.User
		}{
			Title: "Home",
			Year:  time.Now().Year(),
			User:  user, // now set if logged in
		}
		err := templates.ExecuteTemplate(w, "layout", data)
		if err != nil {
			log.Printf("template execution error: %v", err)
			return
		}
	})

	// Read host and port from environment variables
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "10000"
	}
	addr := host + ":" + port
	log.Printf("Server starting on %s...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
