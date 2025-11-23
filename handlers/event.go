package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

func EventHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var input struct {
			Date        string `json:"date"`
			Time        string `json:"time"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		// Combine date and time
		eventTime, err := time.Parse("2006-01-02 15:04", input.Date+" "+input.Time)
		if err != nil {
			http.Error(w, "Invalid date/time format", http.StatusBadRequest)
			return
		}
		event := models.Event{
			Date:        eventTime,
			Description: input.Description,
		}
		if err := event.Create(); err != nil {
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
