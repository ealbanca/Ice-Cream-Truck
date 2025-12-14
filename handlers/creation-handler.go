package handlers

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// CreationsHandler serves the user's custom creations page
func CreationsHandler(w http.ResponseWriter, r *http.Request) error {
	userID, err := GetSessionUserID(r)
	if err != nil || userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil
	}
	user, _ := models.GetUserByID(userID)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil
	}
	creations, _ := models.GetProductsByUserID(user.ID)

	// Enrich creations with Size, Flavors, Toppings fields
	sizesMap := make(map[int]string)
	if sizes, err := models.GetAllSizes(); err == nil {
		for _, s := range sizes {
			sizesMap[s.ID] = s.Label
		}
	}
	flavorsMap := make(map[int]string)
	if flavors, err := models.GetAllFlavors(); err == nil {
		for _, f := range flavors {
			flavorsMap[f.ID] = f.Name
		}
	}
	toppingsMap := make(map[int]string)
	if toppings, err := models.GetAllToppings(); err == nil {
		for _, t := range toppings {
			toppingsMap[t.ID] = t.Name
		}
	}
	type CreationView struct {
		ID          int
		Name        string
		Size        string
		Flavors     string
		Toppings    string
		DateCreated string // Placeholder, add if available in DB
	}
	var creationViews []CreationView
	for _, c := range creations {
		var flavorNames []string
		if c.FlavorID1.Valid {
			flavorNames = append(flavorNames, flavorsMap[int(c.FlavorID1.Int64)])
		}
		if c.FlavorID2.Valid {
			flavorNames = append(flavorNames, flavorsMap[int(c.FlavorID2.Int64)])
		}
		if c.FlavorID3.Valid {
			flavorNames = append(flavorNames, flavorsMap[int(c.FlavorID3.Int64)])
		}
		var toppingNames []string
		if c.ToppingID1.Valid {
			toppingNames = append(toppingNames, toppingsMap[int(c.ToppingID1.Int64)])
		}
		if c.ToppingID2.Valid {
			toppingNames = append(toppingNames, toppingsMap[int(c.ToppingID2.Int64)])
		}
		if c.ToppingID3.Valid {
			toppingNames = append(toppingNames, toppingsMap[int(c.ToppingID3.Int64)])
		}
		creationViews = append(creationViews, CreationView{
			ID:          c.ID,
			Name:        c.ProductName,
			Size:        sizesMap[c.SizeID],
			Flavors:     joinNonEmpty(flavorNames, ", "),
			Toppings:    joinNonEmpty(toppingNames, ", "),
			DateCreated: "", // Not available unless added to DB
		})
	}

	tmpl := template.Must(template.ParseFiles(
		"views/creation/creation-list.gohtml",
		"views/layouts/layout.gohtml",
		"views/partials/head.gohtml",
		"views/partials/header.gohtml",
		"views/partials/navigation.gohtml",
		"views/partials/footer.gohtml",
	))
	data := struct {
		Title     string
		Year      int
		User      *models.User
		Creations []CreationView
	}{
		Title:     "Your Creations",
		Year:      time.Now().Year(),
		User:      user,
		Creations: creationViews,
	}
	return tmpl.ExecuteTemplate(w, "layout", data)
}

// joinNonEmpty joins non-empty strings with sep
func joinNonEmpty(items []string, sep string) string {
	out := []string{}
	for _, s := range items {
		if s != "" {
			out = append(out, s)
		}
	}
	return strings.Join(out, sep)
}
