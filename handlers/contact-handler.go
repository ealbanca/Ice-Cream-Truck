package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// Handles GET (renders page) and POST (saves contact form)
func ContactHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var input struct {
			Name    string `json:"name"`
			Email   string `json:"email"`
			Phone   string `json:"phone"`
			Reason  string `json:"reason"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		contact := models.Contact{
			Name:    input.Name,
			Email:   input.Email,
			Phone:   input.Phone,
			Reason:  input.Reason,
			Message: input.Message,
		}
		if err := contact.Create(); err != nil {
			http.Error(w, "Failed to save contact", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(contact)
	case http.MethodGet:
		contacts, err := models.GetAllContacts()
		if err != nil {
			http.Error(w, "Failed to fetch contacts", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(contacts)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
