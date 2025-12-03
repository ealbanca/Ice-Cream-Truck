package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// RegisterHandler handles GET (renders form) and POST (registers user)
func RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	var message string
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			message = "Invalid form data."
		} else {
			user := models.User{
				Name:     r.FormValue("name"),
				Lastname: r.FormValue("lastname"),
				Username: r.FormValue("username"),
				Email:    r.FormValue("email"),
				Phone:    r.FormValue("phone"),
				Password: r.FormValue("password"),
			}
			if user.Password != r.FormValue("confirmPassword") {
				message = "Passwords do not match."
			} else {
				err = user.Register()
				if err != nil {
					message = "Registration failed: " + err.Error()
				} else {
					message = "Registration successful! You can now log in."
				}
			}
		}
	}

	tmpl := template.Must(template.ParseFiles(
		filepath.Join("views", "layouts", "layout.gohtml"),
		filepath.Join("views", "partials", "head.gohtml"),
		filepath.Join("views", "partials", "header.gohtml"),
		filepath.Join("views", "partials", "navigation.gohtml"),
		filepath.Join("views", "partials", "footer.gohtml"),
		filepath.Join("views", "account", "register.gohtml"),
	))
	data := struct {
		Title   string
		Year    int
		Message string
	}{
		Title:   "Register",
		Year:    time.Now().Year(),
		Message: message,
	}
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		return fmt.Errorf("template execution error: %w", err)
	}
	return nil
}
