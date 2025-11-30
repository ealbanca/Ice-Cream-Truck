package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

func EventHandler(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		// Parse combined date-time string (e.g., 2025-12-01T18:00)
		eventTime, err := time.Parse("2006-01-02T15:04", input.Date)
		if err != nil {
			http.Error(w, "Invalid date/time format", http.StatusBadRequest)
			return
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
			http.Error(w, "Failed to create event", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(event)
	case http.MethodGet:
		events, err := models.GetAllEvents()
		if err != nil {
			http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(events)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
