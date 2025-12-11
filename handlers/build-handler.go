package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

func BuildHandler(w http.ResponseWriter, r *http.Request) error {
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
		filepath.Join("views", "build", "build.gohtml"),
	))
	data := struct {
		Title   string
		Year    int
		User    *models.User
		Message string
	}{
		Title:   "Build Your Own Ice Cream",
		Year:    time.Now().Year(),
		User:    user,
		Message: "",
	}
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		return err
	}
	return nil
}
