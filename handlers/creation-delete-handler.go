package handlers

import (
	"net/http"
	"strconv"

	"github.com/ealbanca/Ice-Cream-Truck/models"
)

// CreationDeleteHandler deletes a custom product by ID
func CreationDeleteHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}
	idStr := r.URL.Path[len("/creation/") : len(r.URL.Path)-len("/delete")]
	idStr = trimSlashes(idStr)
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return nil
	}
	err = models.DeleteCustomProduct(productID)
	if err != nil {
		http.Error(w, "Failed to delete creation", http.StatusInternalServerError)
		return nil
	}
	http.Redirect(w, r, "/account/creations", http.StatusSeeOther)
	return nil
}

func trimSlashes(s string) string {
	for len(s) > 0 && (s[0] == '/' || s[len(s)-1] == '/') {
		if s[0] == '/' {
			s = s[1:]
		}
		if len(s) > 0 && s[len(s)-1] == '/' {
			s = s[:len(s)-1]
		}
	}
	return s
}
