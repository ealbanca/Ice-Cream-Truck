package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ealbanca/Ice-Cream-Truck/config"
)

// CleanupAbandonedGuestCarts deletes guest cart data (orders, order_items, guest_carts, products) older than 1 day
func CleanupAbandonedGuestCarts() error {
	rows, err := config.DB.Query(`
		SELECT id, guest_cart_id FROM orders
		WHERE status = 'guest_cart' AND order_date < NOW() - INTERVAL '1 day'
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var orderIDs []int
	var guestCartIDs []string
	for rows.Next() {
		var id int
		var guestCartID string
		if err := rows.Scan(&id, &guestCartID); err == nil {
			orderIDs = append(orderIDs, id)
			guestCartIDs = append(guestCartIDs, guestCartID)
		}
	}

	// Delete order_items for these orders
	for _, orderID := range orderIDs {
		config.DB.Exec("DELETE FROM order_items WHERE order_id = $1", orderID)
	}

	// Delete guest_carts for these guest_cart_ids
	for _, guestCartID := range guestCartIDs {
		config.DB.Exec("DELETE FROM guest_carts WHERE guest_cart_id = $1", guestCartID)
	}

	// Delete the guest orders themselves
	for _, orderID := range orderIDs {
		config.DB.Exec("DELETE FROM orders WHERE id = $1", orderID)
	}

	// Optionally, delete products that are not referenced by any order_items or guest_carts
	config.DB.Exec(`
		DELETE FROM products
		WHERE id NOT IN (SELECT product_id FROM order_items)
		  AND id NOT IN (SELECT product_id FROM guest_carts)
	`)

	return nil
}

// Add product to guest_carts for a guest_cart_id
func AddProductToGuestCart(guestCartID string, productID int) error {
	if guestCartID == "" {
		// Log and fail if guest_cart_id is empty
		// fmt.Println("Guest cart add failed: guest_cart_id is empty")
		return fmt.Errorf("guest_cart_id is empty")
	}
	var existingID int
	err := config.DB.QueryRow("SELECT id FROM guest_carts WHERE guest_cart_id = $1 AND product_id = $2", guestCartID, productID).Scan(&existingID)
	switch {
	case err == nil:
		// Exists, increment quantity
		_, err = config.DB.Exec("UPDATE guest_carts SET quantity = quantity + 1 WHERE id = $1", existingID)
		if err != nil {
			// fmt.Println("Guest cart update error:", err)
			return err
		}
	case err == sql.ErrNoRows:
		// Not found, insert new
		_, err = config.DB.Exec("INSERT INTO guest_carts (guest_cart_id, product_id, quantity) VALUES ($1, $2, 1)", guestCartID, productID)
		if err != nil {
			// fmt.Println("Guest cart insert error:", err)
			return err
		}
	default:
		// Some other error
		// fmt.Println("Guest cart select error:", err)
		return err
	}

	// --- Insert into order_items for guest order ---
	// Get or create a guest order (user_id NULL, status = 'guest_cart')
	var orderID int
	oerr := config.DB.QueryRow("SELECT id FROM orders WHERE user_id IS NULL AND status = 'guest_cart' AND guest_cart_id = $1", guestCartID).Scan(&orderID)
	if oerr == sql.ErrNoRows {
		// Create new guest order
		oerr = config.DB.QueryRow("INSERT INTO orders (user_id, total_price, status, guest_cart_id) VALUES (NULL, 0, 'guest_cart', $1) RETURNING id", guestCartID).Scan(&orderID)
	}
	if oerr != nil {
		// Log but do not fail guest cart add if order creation fails
		// fmt.Println("Guest order creation error:", oerr)
		return nil
	}
	// Insert product into order_items for this guest order
	_, _ = config.DB.Exec("INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, 1)", orderID, productID)
	// Ignore order_items error for guest cart
	return nil
}

// Update the quantity of a guest cart item
func UpdateGuestCartItemQuantity(itemID int, quantity int) error {
	if quantity < 1 {
		quantity = 1
	}
	// Update guest_carts
	_, err := config.DB.Exec("UPDATE guest_carts SET quantity = $1 WHERE id = $2", quantity, itemID)
	if err != nil {
		return err
	}
	// Also update order_items for the guest order
	// Find guest_cart_id and product_id for this item
	var guestCartID string
	var productID int
	q := "SELECT guest_cart_id, product_id FROM guest_carts WHERE id = $1"
	err = config.DB.QueryRow(q, itemID).Scan(&guestCartID, &productID)
	if err != nil {
		return nil // guest_carts updated, but can't update order_items
	}
	// Find guest order
	var orderID int
	oerr := config.DB.QueryRow("SELECT id FROM orders WHERE user_id IS NULL AND status = 'guest_cart' AND guest_cart_id = $1", guestCartID).Scan(&orderID)
	if oerr != nil {
		return nil // guest_carts updated, but no guest order
	}
	// Update order_items
	config.DB.Exec("UPDATE order_items SET quantity = $1 WHERE order_id = $2 AND product_id = $3", quantity, orderID, productID)
	return nil
}

// Remove a product from the guest cart
func RemoveFromGuestCart(itemID int) error {
	// Find guest_cart_id and product_id for this item
	var guestCartID string
	var productID int
	q := "SELECT guest_cart_id, product_id FROM guest_carts WHERE id = $1"
	err := config.DB.QueryRow(q, itemID).Scan(&guestCartID, &productID)
	if err != nil {
		// Still try to delete from guest_carts
		_, err2 := config.DB.Exec("DELETE FROM guest_carts WHERE id = $1", itemID)
		return err2
	}
	// Find guest order
	var orderID int
	oerr := config.DB.QueryRow("SELECT id FROM orders WHERE user_id IS NULL AND status = 'guest_cart' AND guest_cart_id = $1", guestCartID).Scan(&orderID)
	// Remove from order_items if order found
	if oerr == nil {
		config.DB.Exec("DELETE FROM order_items WHERE order_id = $1 AND product_id = $2", orderID, productID)
	}
	// Remove from guest_carts
	_, err = config.DB.Exec("DELETE FROM guest_carts WHERE id = $1", itemID)
	return err
}

// Get all cart items for a guest_cart_id
type GuestCartItem struct {
	ID            int
	ProductName   string
	SizeLabel     string
	FlavorLabels  string
	ToppingLabels string
	Quantity      int
	TotalPrice    float64
}

func GetGuestCartItems(guestCartID string) []GuestCartItem {
	rows, err := config.DB.Query(`
	   SELECT gc.id, p.product_name, s.size_label,
		  f1.flavor_name, f2.flavor_name, f3.flavor_name,
		  t1.topping_name, t2.topping_name, t3.topping_name,
		  gc.quantity, p.total_price
	   FROM guest_carts gc
	   JOIN products p ON gc.product_id = p.id
	   JOIN sizes s ON p.size_id = s.id
	   LEFT JOIN flavors f1 ON p.flavor_id1 = f1.id
	   LEFT JOIN flavors f2 ON p.flavor_id2 = f2.id
	   LEFT JOIN flavors f3 ON p.flavor_id3 = f3.id
	   LEFT JOIN toppings t1 ON p.topping_id1 = t1.id
	   LEFT JOIN toppings t2 ON p.topping_id2 = t2.id
	   LEFT JOIN toppings t3 ON p.topping_id3 = t3.id
	   WHERE gc.guest_cart_id = $1
	   ORDER BY gc.id DESC
       `, guestCartID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var items []GuestCartItem
	for rows.Next() {
		var c GuestCartItem
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
		c.TotalPrice = c.TotalPrice * float64(c.Quantity)
		items = append(items, c)
	}
	return items
}

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

func SaveCustomProduct(p CustomProduct) (int, error) {
	// Insert new product and return its ID
	var productID int
	err := config.DB.QueryRow(
		`INSERT INTO products (product_name, size_id, flavor_id1, flavor_id2, flavor_id3, topping_id1, topping_id2, topping_id3, total_price)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		p.ProductName, p.SizeID,
		p.FlavorID1, p.FlavorID2, p.FlavorID3,
		p.ToppingID1, p.ToppingID2, p.ToppingID3,
		p.TotalPrice,
	).Scan(&productID)
	return productID, err
}

// Get or create a cart order for the user (status = 'cart')
func GetOrCreateCartOrderID(userID int) (int, error) {
	var orderID int
	err := config.DB.QueryRow("SELECT id FROM orders WHERE user_id = $1 AND status = 'cart'", userID).Scan(&orderID)
	if err == sql.ErrNoRows {
		// Create new cart order
		err = config.DB.QueryRow("INSERT INTO orders (user_id, total_price, status) VALUES ($1, 0, 'cart') RETURNING id", userID).Scan(&orderID)
	}
	return orderID, err
}

// Add product to order_items for the given order
func AddProductToOrderItems(orderID, productID int) error {
	_, err := config.DB.Exec("INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, 1)", orderID, productID)
	return err
}

// GetProductsByUserID returns all products created by a specific user
func GetProductsByUserID(userID int) ([]CustomProduct, error) {
	rows, err := config.DB.Query(`
		SELECT id, product_name, size_id, flavor_id1, flavor_id2, flavor_id3, topping_id1, topping_id2, topping_id3, total_price
		FROM products
		WHERE id IN (
			SELECT product_id FROM order_items oi
			JOIN orders o ON oi.order_id = o.id
			WHERE o.user_id = $1
		)
		ORDER BY id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []CustomProduct
	for rows.Next() {
		var p CustomProduct
		err := rows.Scan(&p.ID, &p.ProductName, &p.SizeID, &p.FlavorID1, &p.FlavorID2, &p.FlavorID3, &p.ToppingID1, &p.ToppingID2, &p.ToppingID3, &p.TotalPrice)
		if err == nil {
			products = append(products, p)
		}
	}
	return products, nil
}

// DeleteCustomProduct deletes a product from the products table by ID
func DeleteCustomProduct(productID int) error {
	// Remove from order_items first to avoid FK constraint errors
	_, _ = config.DB.Exec("DELETE FROM order_items WHERE product_id = $1", productID)
	// Remove from guest_carts as well
	_, _ = config.DB.Exec("DELETE FROM guest_carts WHERE product_id = $1", productID)
	// Now delete the product
	_, err := config.DB.Exec("DELETE FROM products WHERE id = $1", productID)
	return err
}
