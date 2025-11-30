package models

import (
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/config"
)

type Contact struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Reason    string    `json:"reason"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// Create a new contact entry
func (c *Contact) Create() error {
	loc, _ := time.LoadLocation("America/Denver")
	nowDenver := time.Now().In(loc)
	query := `INSERT INTO contacts (name, email, phone, reason, message, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	return config.DB.QueryRow(query, c.Name, c.Email, c.Phone, c.Reason, c.Message, nowDenver).Scan(&c.ID, &c.CreatedAt)
}

// Get all contact entries
func GetAllContacts() ([]Contact, error) {
	rows, err := config.DB.Query("SELECT id, name, email, phone, reason, message, created_at FROM contacts ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contacts []Contact
	for rows.Next() {
		var c Contact
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Reason, &c.Message, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}
