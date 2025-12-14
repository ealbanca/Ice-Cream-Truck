package handlers

import (
	"database/sql"
	"fmt"
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
	var guestCartID string
	var userProducts []models.CustomProduct
	if userID, err := GetSessionUserID(r); err == nil {
		user, _ = models.GetUserByID(userID)
		if user != nil {
			userProducts, _ = models.GetProductsByUserID(user.ID)
		}
	} else {
		// Guest: check for guest_cart_id cookie, or set one
		c, err := r.Cookie("guest_cart_id")
		if err == nil {
			guestCartID = c.Value
		} else {
			// Generate a random guest cart ID (timestamp-based for simplicity)
			guestCartID = strconv.FormatInt(time.Now().UnixNano(), 36)
			http.SetCookie(w, &http.Cookie{
				Name:     "guest_cart_id",
				Value:    guestCartID,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   60 * 60 * 24 * 7, // 1 week
			})
		}
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
			// Flavors have no price, so skip price calculation for flavors
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
			productID, err := models.SaveCustomProduct(product)
			if err != nil {
				message = "Failed to save your custom ice cream."
				messageType = "error"
			} else if user != nil {
				// Add to user's cart order
				orderID, oerr := models.GetOrCreateCartOrderID(user.ID)
				if oerr == nil {
					piErr := models.AddProductToOrderItems(orderID, productID)
					if piErr == nil {
						message = "Your custom ice cream has been added to your cart!"
						messageType = "success"
					} else {
						message = "Product saved, but failed to add to cart."
						messageType = "error"
					}
				} else {
					message = "Product saved, but failed to create cart order."
					messageType = "error"
				}
			} else if guestCartID != "" {
				// Guest cart: add to guest_carts table
				piErr := models.AddProductToGuestCart(guestCartID, productID)
				if piErr == nil {
					message = "Your custom ice cream has been added to your guest cart!"
					messageType = "success"
				} else {
					// Log the error for debugging
					fmt.Println("Guest cart add error:", piErr)
					message = "Product saved, but failed to add to guest cart."
					messageType = "error"
				}
			} else {
				message = "Product saved, but could not add to cart."
				messageType = "error"
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
		Title        string
		Year         int
		User         *models.User
		Sizes        []models.Size
		Flavors      []models.Flavor
		Toppings     []models.Topping
		UserProducts []models.CustomProduct
		Message      string
		MessageType  string
	}{
		Title:        "Build Your Own Ice Cream",
		Year:         time.Now().Year(),
		User:         user,
		Sizes:        sizes,
		Flavors:      flavors,
		Toppings:     toppings,
		UserProducts: userProducts,
		Message:      message,
		MessageType:  messageType,
	}
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		return err
	}
	return nil
}
