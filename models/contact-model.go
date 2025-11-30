package models

import (
	"github.com/ealbanca/Ice-Cream-Truck/config"
)

type Contact struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// Create a new contact entry
func (c *Contact) Create() error {
	query := `INSERT INTO contacts (name, email, phone, reason, message) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return config.DB.QueryRow(query, c.Name, c.Email, c.Phone, c.Reason, c.Message).Scan(&c.ID)
}

// Get all contact entries
func GetAllContacts() ([]Contact, error) {
	rows, err := config.DB.Query("SELECT id, name, email, phone, reason, message FROM contacts ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contacts []Contact
	for rows.Next() {
		var c Contact
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Reason, &c.Message)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}
