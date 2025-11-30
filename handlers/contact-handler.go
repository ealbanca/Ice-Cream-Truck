package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// Handles GET (renders page) and POST (saves contact form)
func ContactHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var input struct {
			ID        int       `json:"id"`
			Name      string    `json:"name"`
			Email     string    `json:"email"`
			Phone     string    `json:"phone"`
			Reason    string    `json:"reason"`
			Message   string    `json:"message"`
			CreatedAt time.Time `json:"created_at"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		contact := models.Contact{
			Name:    input.Name,
			Email:   input.Email,
			Phone:   input.Phone,
			Reason:  input.Reason,
			Message: input.Message,
		}
		if err := contact.Create(); err != nil {
			http.Error(w, "Failed to save contact", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(contact)
	case http.MethodGet:
		// Render the contact page template
		tmpl := template.Must(template.ParseFiles(
			filepath.Join("views", "layouts", "layout.gohtml"),
			filepath.Join("views", "partials", "head.gohtml"),
			filepath.Join("views", "partials", "header.gohtml"),
			filepath.Join("views", "partials", "navigation.gohtml"),
			filepath.Join("views", "partials", "footer.gohtml"),
			filepath.Join("views", "contact", "contact.gohtml"),
		))
		data := struct {
			Title string
			Year  int
		}{
			Title: "Contact",
			Year:  time.Now().Year(),
		}
		err := tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
