package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("views", "layouts", "layout.gohtml"),
		filepath.Join("views", "partials", "head.gohtml"),
		filepath.Join("views", "partials", "header.gohtml"),
		filepath.Join("views", "partials", "navigation.gohtml"),
		filepath.Join("views", "partials", "footer.gohtml"),
		filepath.Join("views", "account", "login.gohtml"),
	))
	data := struct {
		Title string
		Year  int
	}{
		Title: "Login",
		Year:  time.Now().Year(),
	}
	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
