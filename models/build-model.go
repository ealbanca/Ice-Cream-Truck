package models

import (
	"database/sql"

	"github.com/ealbanca/Ice-Cream-Truck/config"
)

type Size struct {
	ID    int
	Label string
	Price float64
	Img   string
}

type Flavor struct {
	ID   int
	Name string
	Img  string
}

type Topping struct {
	ID    int
	Name  string
	Price float64
	Img   string
}

type CustomProduct struct {
	ID          int
	ProductName string
	UserID      int
	SizeID      int
	FlavorID1   sql.NullInt64
	FlavorID2   sql.NullInt64
	FlavorID3   sql.NullInt64
	ToppingID1  sql.NullInt64
	ToppingID2  sql.NullInt64
	ToppingID3  sql.NullInt64
	Quantity    int
	TotalPrice  float64
}

func GetAllSizes() ([]Size, error) {
	rows, err := config.DB.Query("SELECT id, size_label, price, size_img FROM sizes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sizes []Size
	for rows.Next() {
		var s Size
		if err := rows.Scan(&s.ID, &s.Label, &s.Price, &s.Img); err == nil {
			sizes = append(sizes, s)
		}
	}
	return sizes, nil
}

func GetAllFlavors() ([]Flavor, error) {
	rows, err := config.DB.Query("SELECT id, flavor_name, falvor_img FROM flavors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var flavors []Flavor
	for rows.Next() {
		var f Flavor
		if err := rows.Scan(&f.ID, &f.Name, &f.Img); err == nil {
			flavors = append(flavors, f)
		}
	}
	return flavors, nil
}

func GetAllToppings() ([]Topping, error) {
	rows, err := config.DB.Query("SELECT id, topping_name, additional_price, topping_img FROM toppings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var toppings []Topping
	for rows.Next() {
		var t Topping
		if err := rows.Scan(&t.ID, &t.Name, &t.Price, &t.Img); err == nil {
			toppings = append(toppings, t)
		}
	}
	return toppings, nil
}

func SaveCustomProduct(p CustomProduct) error {
	// Check for existing identical product for this user (robust null handling)
	var existingID int
	err := config.DB.QueryRow(`
		SELECT id FROM products WHERE user_id = $1 AND product_name = $2 AND size_id = $3
		AND flavor_id1 IS NOT DISTINCT FROM $4
		AND flavor_id2 IS NOT DISTINCT FROM $5
		AND flavor_id3 IS NOT DISTINCT FROM $6
		AND topping_id1 IS NOT DISTINCT FROM $7
		AND topping_id2 IS NOT DISTINCT FROM $8
		AND topping_id3 IS NOT DISTINCT FROM $9
	`,
		p.UserID, p.ProductName, p.SizeID,
		p.FlavorID1, p.FlavorID2, p.FlavorID3,
		p.ToppingID1, p.ToppingID2, p.ToppingID3,
	).Scan(&existingID)
	if err == nil {
		// Product exists, increment quantity
		_, err = config.DB.Exec(`UPDATE products SET quantity = quantity + 1 WHERE id = $1`, existingID)
		return err
	}
	// If not found, insert new
	_, err = config.DB.Exec(
		`INSERT INTO products (product_name, user_id, size_id, flavor_id1, flavor_id2, flavor_id3, topping_id1, topping_id2, topping_id3, quantity, total_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		p.ProductName, p.UserID, p.SizeID,
		p.FlavorID1, p.FlavorID2, p.FlavorID3,
		p.ToppingID1, p.ToppingID2, p.ToppingID3,
		1, p.TotalPrice,
	)
	return err
}
