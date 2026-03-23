package osm

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dennstack/addrgo/db"
)

func ImportOSMData() {
	urls := os.Getenv("OSM_URLS")
	if urls == "" {
		fmt.Println("OSM_URLS not set, skipping import.")
		return
	}
	urlList := strings.Split(urls, ",")
	db := db.GetDatabaseInstance().Connection

	fmt.Println("Starting OSM data import...")
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS addrgo_addresses (street VARCHAR(255), city VARCHAR(100), postcode VARCHAR(20), hash VARCHAR(16), UNIQUE KEY uk_hash (hash), INDEX idx_postcode (postcode), INDEX idx_street_city (street, city), INDEX idx_city (city), INDEX idx_street (street), INDEX idx_street_postcode (street, postcode), INDEX idx_full_address (street, city, postcode))")
	if err != nil {
		fmt.Println("Error creating addrgo_addresses table:", err)
		return
	}

	stream := make(chan Address)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		saveAddress(db, stream)
	}()

	for _, url := range urlList {
		fmt.Printf("Processing URL: %s\n", url)
		ParseFromUrl(url, stream)
	}
	close(stream)
	wg.Wait()
	fmt.Println("OSM import complete.")
}

func saveAddress(db *sql.DB, stream <-chan Address) {
	for address := range stream {
		_, err := db.Exec(
			"INSERT IGNORE INTO addrgo_addresses (street, city, postcode, hash) VALUES (?, ?, ?, ?)",
			address.Street, address.City, address.Postcode, address.Hash,
		)
		if err != nil {
			fmt.Println("Error inserting address:", err)
		}
	}
}
