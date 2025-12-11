package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// Handles GET (renders page) and POST (saves contact form)
func ContactHandler(w http.ResponseWriter, r *http.Request) error {
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
			return fmt.Errorf("invalid input: %w", err)
		}
		contact := models.Contact{
			Name:    input.Name,
			Email:   input.Email,
			Phone:   input.Phone,
			Reason:  input.Reason,
			Message: input.Message,
		}
		if err := contact.Create(); err != nil {
			return fmt.Errorf("failed to save contact: %w", err)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(contact)
		return nil
	case http.MethodGet:
		// Render the contact page template
		var user *models.User
		if userID, err := GetSessionUserID(r); err == nil {
			user, _ = models.GetUserByID(userID)
		}
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
			User  *models.User
		}{
			Title: "Contact",
			Year:  time.Now().Year(),
			User:  user,
		}
		if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("method not allowed")
	}
}
