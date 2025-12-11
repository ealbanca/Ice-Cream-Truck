package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/config"
	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// Helper: Set session cookie with user ID
func setSession(w http.ResponseWriter, userID int) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    strconv.Itoa(userID),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

// Helper: Get user ID from session cookie
func GetSessionUserID(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(cookie.Value)
}

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
			} else {
				setSession(w, user.ID)
			}
		}
	}

	// If user is logged in, show management page
	if user == nil {
		// Try to load user from session
		if userID, err := GetSessionUserID(r); err == nil {
			user, _ = models.GetUserByID(userID)
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

	// Show login page if not logged in
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
		User    *models.User
	}{
		Title:   "Login",
		Year:    time.Now().Year(),
		Message: message,
		User:    user, // will be nil if not logged in
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
			name := r.FormValue("name")
			lastname := r.FormValue("lastname")
			username := r.FormValue("username")
			email := r.FormValue("email")
			phone := r.FormValue("phone")
			password := r.FormValue("password")
			confirmPassword := r.FormValue("confirmPassword")
			if password != confirmPassword {
				message = "Passwords do not match."
			} else {
				user := &models.User{
					Name:      name,
					Lastname:  lastname,
					Username:  username,
					Email:     email,
					Phone:     phone,
					Password:  password, // In production, hash the password!
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

// LogoutHandler logs the user out by clearing the session/cookie and redirecting to login
func LogoutHandler(w http.ResponseWriter, r *http.Request) error {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Delete cookie
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}

// EditAccountHandler handles GET and POST for editing user account
func EditAccountHandler(w http.ResponseWriter, r *http.Request) error {
	var message string
	var user *models.User

	// Load user from session
	if userID, err := GetSessionUserID(r); err == nil {
		user, _ = models.GetUserByID(userID)
	}
	if user == nil {
		message = "User not found."
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			message = "Invalid form data."
		} else if user != nil {
			// Update user fields from form values
			user.Name = r.FormValue("name")
			user.Lastname = r.FormValue("lastname")
			user.Username = r.FormValue("username")
			user.Email = r.FormValue("email")
			user.Phone = r.FormValue("phone")

			password := r.FormValue("password")
			confirmPassword := r.FormValue("confirm-password")
			if password != "" {
				if password != confirmPassword {
					message = "Passwords do not match."
				} else {
					err := models.UpdateUserPassword(user.ID, password)
					if err != nil {
						message = "Failed to update password: " + err.Error()
					} else {
						message = "Password updated successfully! "
					}
				}
			}

			err := models.UpdateUser(user)
			if err != nil {
				message += "Failed to update account: " + err.Error()
			} else if password == "" || password == confirmPassword {
				message += "Account updated successfully!"
			}
		}
	}
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("views", "layouts", "layout.gohtml"),
		filepath.Join("views", "partials", "head.gohtml"),
		filepath.Join("views", "partials", "header.gohtml"),
		filepath.Join("views", "partials", "navigation.gohtml"),
		filepath.Join("views", "partials", "footer.gohtml"),
		filepath.Join("views", "account", "edit-account.gohtml"),
	))
	data := struct {
		Title   string
		Year    int
		User    *models.User
		Message string
	}{
		Title:   "Edit Account",
		Year:    time.Now().Year(),
		User:    user,
		Message: message,
	}
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		return err
	}
	return nil
}
