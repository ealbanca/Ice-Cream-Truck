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
	SizeID      int
	FlavorID1   sql.NullInt64
	FlavorID2   sql.NullInt64
	FlavorID3   sql.NullInt64
	ToppingID1  sql.NullInt64
	ToppingID2  sql.NullInt64
	ToppingID3  sql.NullInt64
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
	_, err := config.DB.Exec(
		`INSERT INTO products (product_name, size_id, flavor_id1, flavor_id2, flavor_id3, topping_id1, topping_id2, topping_id3, total_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		p.ProductName, p.SizeID, p.FlavorID1, p.FlavorID2, p.FlavorID3, p.ToppingID1, p.ToppingID2, p.ToppingID3, p.TotalPrice,
	)
	return err
}
