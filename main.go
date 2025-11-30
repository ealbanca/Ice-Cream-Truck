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
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file (for local development)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file (this is OK in production)")
	}

	// Initialize database connection
	config.InitDB()
	defer config.DB.Close()

	// Initialize routes
	http.HandleFunc("/api/products", handlers.ProductsHandler)
	http.HandleFunc("/api/events", handlers.EventHandler)
	http.HandleFunc("/about", handlers.AboutHandler)
	http.HandleFunc("/contact", handlers.ContactHandler)
	http.HandleFunc("/api/contact", handlers.ContactHandler)

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
		data := struct {
			Title string
			Year  int
		}{
			Title: "Home",
			Year:  time.Now().Year(),
		}
		err := templates.ExecuteTemplate(w, "layout", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
