package models

import (
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/config"
)

type Event struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Date        time.Time `json:"date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Create a new event
func (e *Event) Create() error {
	e.CreatedAt = time.Now()
	query := `INSERT INTO events (name, email, phone, date, start_time, end_time, description, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return config.DB.QueryRow(query, e.Name, e.Email, e.Phone, e.Date, e.StartTime, e.EndTime, e.Description, e.CreatedAt).Scan(&e.ID)
}

// Get all events
func GetAllEvents() ([]Event, error) {
	rows, err := config.DB.Query("SELECT id, name, email, phone, date, start_time, end_time, description, created_at FROM events ORDER BY date")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.Phone, &e.Date, &e.StartTime, &e.EndTime, &e.Description, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
