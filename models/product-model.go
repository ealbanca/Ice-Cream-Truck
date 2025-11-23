package models

import (
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/config"
)

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) Create() error {
	query := `INSERT INTO products (name, description, price, created_at, updated_at) 
             VALUES (?, ?, ?, NOW(), NOW())`

	result, err := config.DB.Exec(query, p.Name, p.Description, p.Price)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int(id)
	return nil
}

func GetAllProducts() ([]Product, error) {
	var products []Product

	rows, err := config.DB.Query("SELECT id, name, description, price, created_at, updated_at FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
