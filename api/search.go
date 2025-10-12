package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dennstack/addrgo/db"
)

type SearchRequest struct {
	Limit    int    `json:"limit,omitempty"`
	Street   string `json:"street,omitempty"`
	City     string `json:"city,omitempty"`
	Postcode string `json:"postcode,omitempty"`
}

type SearchResult struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	Postcode string `json:"postcode"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Count   int            `json:"count"`
}

func buildSearchQuery(request SearchRequest) (string, []interface{}) {
	baseQuery := "SELECT DISTINCT street, city, postcode FROM addrgo_addresses WHERE "
	var conditions []string
	var args []interface{}

	if request.Street != "" {
		conditions = append(conditions, "street LIKE ?")
		args = append(args, "%"+request.Street+"%")
	}
	if request.City != "" {
		conditions = append(conditions, "city LIKE ?")
		args = append(args, "%"+request.City+"%")
	}
	if request.Postcode != "" {
		conditions = append(conditions, "postcode LIKE ?")
		args = append(args, "%"+request.Postcode+"%")
	}

	query := baseQuery + strings.Join(conditions, " AND ")
	query += " ORDER BY postcode, city, street"

	if request.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, request.Limit)
	}

	return query, args
}

func SearcHandler(w http.ResponseWriter, r *http.Request) {
	var request SearchRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	database := db.GetDatabaseInstance().Connection
	query, args := buildSearchQuery(request)

	rows, err := database.Query(query, args...)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var result SearchResult
		err := rows.Scan(&result.Street, &result.City, &result.Postcode)
		if err != nil {
			continue
		}
		results = append(results, result)
	}

	response := SearchResponse{
		Results: results,
		Count:   len(results),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
