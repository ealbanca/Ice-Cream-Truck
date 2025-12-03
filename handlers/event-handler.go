package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

func EventHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
		var input struct {
			ID          int       `json:"id"`
			Name        string    `json:"name"`
			Email       string    `json:"email"`
			Phone       string    `json:"phone"`
			Date        string    `json:"date"`
			Description string    `json:"description"`
			CreatedAt   time.Time `json:"created_at"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			return fmt.Errorf("invalid input: %w", err)
		}
		eventTime, err := time.Parse("2006-01-02T15:04", input.Date)
		if err != nil {
			return fmt.Errorf("invalid date/time format: %w", err)
		}
		event := models.Event{
			Name:        input.Name,
			Email:       input.Email,
			Phone:       input.Phone,
			Date:        eventTime,
			Description: input.Description,
		}
		if err := event.Create(); err != nil {
			log.Println("Event Create error:", err)
			return fmt.Errorf("failed to create event: %w", err)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(event)
		return nil
	case http.MethodGet:
		events, err := models.GetAllEvents()
		if err != nil {
			return fmt.Errorf("failed to fetch events: %w", err)
		}
		json.NewEncoder(w).Encode(events)
		return nil
	default:
		return fmt.Errorf("method not allowed")
	}
}
