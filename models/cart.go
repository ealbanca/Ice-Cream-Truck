package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ealbanca/Ice-Cream-Truck/config"
)

// UpdateCartItemQuantity updates the quantity of a cart item (order_items row)
func UpdateCartItemQuantity(itemID int, quantity int) {
	if quantity < 1 {
		quantity = 1
	}
	config.DB.Exec("UPDATE order_items SET quantity = $1 WHERE id = $2", quantity, itemID)
}

type CartItem struct {
	ID            int
	ProductName   string
	SizeLabel     string
	FlavorLabels  string
	ToppingLabels string
	Quantity      int
	TotalPrice    float64
}

func GetCartItems(userID int) []CartItem {
	if userID == 0 {
		return nil
	}
	// Find the user's active cart order
	var orderID int
	err := config.DB.QueryRow("SELECT id FROM orders WHERE user_id = $1 AND status = 'cart'", userID).Scan(&orderID)
	if err != nil {
		return nil
	}
	rows, err := config.DB.Query(`
	  SELECT oi.id, p.product_name, s.size_label, 
		 f1.flavor_name, f2.flavor_name, f3.flavor_name,
		 t1.topping_name, t2.topping_name, t3.topping_name,
		 oi.quantity, p.total_price
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
	  ORDER BY oi.id DESC
      `, orderID)
	if err != nil {
		fmt.Println("GetCartItems DB error:", err)
		return nil
	}
	defer rows.Close()
	var items []CartItem
	for rows.Next() {
		var c CartItem
		var flavor1, flavor2, flavor3 sql.NullString
		var topping1, topping2, topping3 sql.NullString
		err := rows.Scan(&c.ID, &c.ProductName, &c.SizeLabel, &flavor1, &flavor2, &flavor3, &topping1, &topping2, &topping3, &c.Quantity, &c.TotalPrice)
		if err != nil {
			continue
		}
		// Combine flavor/topping labels
		var flavors, toppings []string
		if flavor1.Valid && flavor1.String != "" {
			flavors = append(flavors, flavor1.String)
		}
		if flavor2.Valid && flavor2.String != "" {
			flavors = append(flavors, flavor2.String)
		}
		if flavor3.Valid && flavor3.String != "" {
			flavors = append(flavors, flavor3.String)
		}
		if topping1.Valid && topping1.String != "" {
			toppings = append(toppings, topping1.String)
		}
		if topping2.Valid && topping2.String != "" {
			toppings = append(toppings, topping2.String)
		}
		if topping3.Valid && topping3.String != "" {
			toppings = append(toppings, topping3.String)
		}
		c.FlavorLabels = strings.Join(flavors, ", ")
		c.ToppingLabels = strings.Join(toppings, ", ")
		// Multiply price by quantity for display
		c.TotalPrice = c.TotalPrice * float64(c.Quantity)
		items = append(items, c)
	}
	return items
}

// RemoveFromCart deletes an item from the user's cart (order_items table)
func RemoveFromCart(itemID int) {
	config.DB.Exec("DELETE FROM order_items WHERE id = $1", itemID)
}

// GetCartTotal calculates the total for the user's active cart order
func GetCartTotal(userID int) float64 {
	var total float64
	var orderID int
	err := config.DB.QueryRow("SELECT id FROM orders WHERE user_id = $1 AND status = 'cart'", userID).Scan(&orderID)
	if err != nil {
		return 0.0
	}
	err = config.DB.QueryRow(`
	       SELECT COALESCE(SUM(p.total_price * oi.quantity),0)
	       FROM order_items oi
	       JOIN products p ON oi.product_id = p.id
	       WHERE oi.order_id = $1
       `, orderID).Scan(&total)
	if err != nil {
		return 0.0
	}
	return total
}
