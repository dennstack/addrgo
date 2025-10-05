package osm

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func ImportOSMData(db *sql.DB) {
	urls := os.Getenv("OSM_URLS")
	urlList := strings.Split(urls, ",")

	fmt.Println("Starting OSM data import...")
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS addrgo_addresses (street VARCHAR(255), city VARCHAR(100), postcode VARCHAR(20), hash VARCHAR(32), INDEX idx_postcode (postcode), INDEX idx_street_city (street, city), INDEX idx_city (city), INDEX idx_street (street), INDEX idx_street_postcode (street, postcode), INDEX idx_full_address (street, city, postcode), INDEX idx_hash (hash))")
	if err != nil {
		fmt.Println("Error creating addrgo_addresses table:", err)
		return
	}

	stream := make(chan Address)
	go saveAddress(db, stream)
	for _, url := range urlList {
		fmt.Printf("Processing URL: %s\n", url)
		ParseFromUrl(url, stream)
	}
	close(stream)

}

func saveAddress(db *sql.DB, stream <-chan Address) error {
	for address := range stream {
		result, err := db.Query("SELECT hash FROM addrgo_addresses WHERE hash = ?", address.Hash)
		if err != nil {
			fmt.Println("Error checking for existing address:", err)
			continue
		}
		if result.Next() {
			// Address already exists, skip insertion
			result.Close()
			continue
		}
		result.Close()
		_, err = db.Exec("INSERT INTO addrgo_addresses (street, city, postcode, hash) VALUES (?, ?, ?, ?)", address.Street, address.City, address.Postcode, address.Hash)
		if err != nil {
			fmt.Println("Error inserting address:", err)
		}
	}
	return nil
}
