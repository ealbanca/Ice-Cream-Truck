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
