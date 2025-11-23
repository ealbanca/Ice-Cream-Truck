package models

import (
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/config"
)

type Event struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

// Create a new event
func (e *Event) Create() error {
	query := `INSERT INTO events (date, description) VALUES ($1, $2) RETURNING id`
	return config.DB.QueryRow(query, e.Date, e.Description).Scan(&e.ID)
}

// Get all events
func GetAllEvents() ([]Event, error) {
	rows, err := config.DB.Query("SELECT id, date, description FROM events ORDER BY date")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Date, &e.Description); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
