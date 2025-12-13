package models

import (
	"github.com/ealbanca/Ice-Cream-Truck/config"
)

type CheckoutInfo struct {
	OrderID    int
	FullName   string
	Email      string
	Phone      string
	CreditCard string
	TotalPrice float64
}

func SaveCheckoutInfo(info CheckoutInfo) error {
	// Insert into checkout_info
	_, err := config.DB.Exec(`
		INSERT INTO checkout_info (order_id, full_name, email, phone, credit_card)
		VALUES ($1, $2, $3, $4, $5)
	`, info.OrderID, info.FullName, info.Email, info.Phone, info.CreditCard)
	if err != nil {
		return err
	}
	// Update orders table: set status to 'completed' and update total_price
	_, err = config.DB.Exec(`
		UPDATE orders SET status = 'completed', total_price = $1 WHERE id = $2
	`, info.TotalPrice, info.OrderID)
	return err
}
