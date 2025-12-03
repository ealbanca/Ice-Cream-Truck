package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) error {
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
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		return fmt.Errorf("template execution error: %w", err)
	}
	return nil
}
