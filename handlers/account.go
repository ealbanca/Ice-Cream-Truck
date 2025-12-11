package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/config"
	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// LoginHandler handles GET and POST for login
func LoginHandler(w http.ResponseWriter, r *http.Request) error {
	var message string
	var user *models.User
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			message = "Invalid form data."
		} else {
			usernameOrEmail := r.FormValue("username")
			password := r.FormValue("password")
			user, err = models.AuthenticateUser(config.DB, usernameOrEmail, password)
			if err != nil {
				message = "Invalid username/email or password."
			}
		}
	}

	if user != nil {
		tmpl := template.Must(template.ParseFiles(
			filepath.Join("views", "layouts", "layout.gohtml"),
			filepath.Join("views", "partials", "head.gohtml"),
			filepath.Join("views", "partials", "header.gohtml"),
			filepath.Join("views", "partials", "navigation.gohtml"),
			filepath.Join("views", "partials", "footer.gohtml"),
			filepath.Join("views", "account", "management.gohtml"),
		))
		data := struct {
			Title string
			Year  int
			User  *models.User
		}{
			Title: "Account Management",
			Year:  time.Now().Year(),
			User:  user,
		}
		if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
			return fmt.Errorf("template execution error: %w", err)
		}
		return nil
	}

	tmpl := template.Must(template.ParseFiles(
		filepath.Join("views", "layouts", "layout.gohtml"),
		filepath.Join("views", "partials", "head.gohtml"),
		filepath.Join("views", "partials", "header.gohtml"),
		filepath.Join("views", "partials", "navigation.gohtml"),
		filepath.Join("views", "partials", "footer.gohtml"),
		filepath.Join("views", "account", "login.gohtml"),
	))
	data := struct {
		Title   string
		Year    int
		Message string
	}{
		Title:   "Login",
		Year:    time.Now().Year(),
		Message: message,
	}
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		return fmt.Errorf("template execution error: %w", err)
	}
	return nil
}

// RegisterHandler handles GET (renders form) and POST (registers user)
func RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	var message string
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			message = "Invalid form data."
		} else {
			username := r.FormValue("username")
			password := r.FormValue("password")
			email := r.FormValue("email")
			// You may want to add validation here
			user := &models.User{
				Username:  username,
				Password:  password, // In production, hash the password!
				Email:     email,
				CreatedAt: time.Now(),
			}
			err := models.RegisterUser(user)
			if err != nil {
				message = "Registration failed: " + err.Error()
			} else {
				message = "Registration successful! Please log in."
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
	data := map[string]interface{}{
		"Message": message,
	}
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		return err
	}
	return nil
}
