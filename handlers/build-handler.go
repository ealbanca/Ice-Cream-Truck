package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/config"
	"github.com/ealbanca/Ice-Cream-Truck/models"
)

func BuildHandler(w http.ResponseWriter, r *http.Request) error {
	var user *models.User
	if userID, err := GetSessionUserID(r); err == nil {
		user, _ = models.GetUserByID(userID)
	}

	message := ""
	messageType := ""

	sizes, _ := models.GetAllSizes()
	flavors, _ := models.GetAllFlavors()
	toppings, _ := models.GetAllToppings()

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			message = "Invalid form data."
			messageType = "error"
		} else {
			sizeID, _ := strconv.Atoi(r.FormValue("size"))
			flavorIDs := strings.Split(r.FormValue("flavor"), ",")
			toppingIDs := strings.Split(r.FormValue("toppings"), ",")

			// Only use up to 3 flavors and 3 toppings (repeats allowed)
			var flavorID1, flavorID2, flavorID3 sql.NullInt64
			if len(flavorIDs) > 0 && flavorIDs[0] != "" {
				id, _ := strconv.Atoi(flavorIDs[0])
				flavorID1 = sql.NullInt64{Int64: int64(id), Valid: true}
			}
			if len(flavorIDs) > 1 {
				id, _ := strconv.Atoi(flavorIDs[1])
				flavorID2 = sql.NullInt64{Int64: int64(id), Valid: true}
			}
			if len(flavorIDs) > 2 {
				id, _ := strconv.Atoi(flavorIDs[2])
				flavorID3 = sql.NullInt64{Int64: int64(id), Valid: true}
			}

			var toppingID1, toppingID2, toppingID3 sql.NullInt64
			if len(toppingIDs) > 0 && toppingIDs[0] != "" {
				id, _ := strconv.Atoi(toppingIDs[0])
				toppingID1 = sql.NullInt64{Int64: int64(id), Valid: true}
			}
			if len(toppingIDs) > 1 {
				id, _ := strconv.Atoi(toppingIDs[1])
				toppingID2 = sql.NullInt64{Int64: int64(id), Valid: true}
			}
			if len(toppingIDs) > 2 {
				id, _ := strconv.Atoi(toppingIDs[2])
				toppingID3 = sql.NullInt64{Int64: int64(id), Valid: true}
			}

			// --- Calculate total price from DB ---
			totalPrice := 0.0
			// Size price
			var sizePrice float64
			err = config.DB.QueryRow("SELECT price FROM sizes WHERE id = $1", sizeID).Scan(&sizePrice)
			if err == nil {
				totalPrice += sizePrice
			}
			// Flavor prices (sum, allow repeats)
			for _, fid := range flavorIDs {
				if fid == "" {
					continue
				}
				var flavorPrice float64
				err = config.DB.QueryRow("SELECT COALESCE(flavor_price, 0) FROM flavor WHERE flavor_id = $1", fid).Scan(&flavorPrice)
				if err == nil {
					totalPrice += flavorPrice
				}
			}
			// Topping prices (sum, allow repeats)
			for _, tid := range toppingIDs {
				if tid == "" {
					continue
				}
				var toppingPrice float64
				err = config.DB.QueryRow("SELECT COALESCE(additional_price, 0) FROM toppings WHERE id = $1", tid).Scan(&toppingPrice)
				if err == nil {
					totalPrice += toppingPrice
				}
			}

			product := models.CustomProduct{
				ProductName: "Custom Ice Cream",
				SizeID:      sizeID,
				FlavorID1:   flavorID1,
				FlavorID2:   flavorID2,
				FlavorID3:   flavorID3,
				ToppingID1:  toppingID1,
				ToppingID2:  toppingID2,
				ToppingID3:  toppingID3,
				TotalPrice:  totalPrice,
			}
			err = models.SaveCustomProduct(product)
			if err != nil {
				message = "Failed to save your custom ice cream."
				messageType = "error"
			} else {
				message = "Your custom ice cream has been saved!"
				messageType = "success"
			}
		}
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
		Title       string
		Year        int
		User        *models.User
		Sizes       []models.Size
		Flavors     []models.Flavor
		Toppings    []models.Topping
		Message     string
		MessageType string
	}{
		Title:       "Build Your Own Ice Cream",
		Year:        time.Now().Year(),
		User:        user,
		Sizes:       sizes,
		Flavors:     flavors,
		Toppings:    toppings,
		Message:     message,
		MessageType: messageType,
	}
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		return err
	}
	return nil
}
