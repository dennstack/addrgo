package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dennstack/addrgo/db"
)

type ValidationRequest struct {
	Street   string `json:"street,omitempty"`
	City     string `json:"city,omitempty"`
	Postcode string `json:"postcode,omitempty"`
}

type ValidationResponse struct {
	Valid bool `json:"valid"`
}

func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	var validate ValidationRequest
	err := json.NewDecoder(r.Body).Decode(&validate)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	database := db.GetDatabaseInstance().Connection
	rows, err := database.Query("SELECT count(*) FROM addrgo_addresses WHERE street=? AND city=? AND postcode=?", validate.Street, validate.City, validate.Postcode)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}

	response := ValidationResponse{
		Valid: count == 1,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
