package models

import (
	"database/sql"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/config"
)

type Order struct {
	ID         int
	UserID     sql.NullInt64
	Date       time.Time
	TotalPrice float64
	Status     string
}

type OrderItem struct {
	Name     string
	Size     string
	Flavors  string
	Toppings string
	Quantity int
	Price    float64
}

// GetOrderByID fetches a single order by its ID
func GetOrderByID(orderID int) (*Order, error) {
	var o Order
	err := config.DB.QueryRow(`SELECT id, user_id, order_date, total_price, status FROM orders WHERE id = $1`, orderID).Scan(&o.ID, &o.UserID, &o.Date, &o.TotalPrice, &o.Status)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// GetOrderItems fetches all items for a given order, including size, flavors, and toppings
func GetOrderItems(orderID int) ([]OrderItem, error) {
	rows, err := config.DB.Query(`
		SELECT 
			p.product_name,
			s.size_label,
			COALESCE(f1.flavor_name, '') ||
				CASE WHEN f2.flavor_name IS NOT NULL THEN ', ' || f2.flavor_name ELSE '' END ||
				CASE WHEN f3.flavor_name IS NOT NULL THEN ', ' || f3.flavor_name ELSE '' END AS flavors,
			TRIM(BOTH ', ' FROM
				COALESCE(t1.topping_name, '') ||
				CASE WHEN t2.topping_name IS NOT NULL THEN ', ' || t2.topping_name ELSE '' END ||
				CASE WHEN t3.topping_name IS NOT NULL THEN ', ' || t3.topping_name ELSE '' END
			) AS toppings,
			oi.quantity,
			p.total_price
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		JOIN sizes s ON p.size_id = s.id
		LEFT JOIN flavors f1 ON p.flavor_id1 = f1.id
		LEFT JOIN flavors f2 ON p.flavor_id2 = f2.id
		LEFT JOIN flavors f3 ON p.flavor_id3 = f3.id
		LEFT JOIN toppings t1 ON p.topping_id1 = t1.id
		LEFT JOIN toppings t2 ON p.topping_id2 = t2.id
		LEFT JOIN toppings t3 ON p.topping_id3 = t3.id
		WHERE oi.order_id = $1
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.Name, &item.Size, &item.Flavors, &item.Toppings, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// GetOrdersByUserID fetches all orders for a given user
func GetOrdersByUserID(userID int) ([]Order, error) {
	rows, err := config.DB.Query(`SELECT id, user_id, order_date, total_price, status FROM orders WHERE user_id = $1 ORDER BY order_date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Date, &o.TotalPrice, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// GetAllOrders fetches all orders in the database (admin view)
func GetAllOrders() ([]Order, error) {
	rows, err := config.DB.Query(`SELECT id, user_id, order_date, total_price, status FROM orders ORDER BY order_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Date, &o.TotalPrice, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
