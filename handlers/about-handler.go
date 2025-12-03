package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) error {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("views", "layouts", "layout.gohtml"),
		filepath.Join("views", "partials", "head.gohtml"),
		filepath.Join("views", "partials", "header.gohtml"),
		filepath.Join("views", "partials", "navigation.gohtml"),
		filepath.Join("views", "partials", "footer.gohtml"),
		filepath.Join("views", "about", "about.gohtml"),
	))
	data := struct {
		Title string
		Year  int
	}{
		Title: "About",
		Year:  time.Now().Year(),
	}
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		return err
	}
	return nil
}
