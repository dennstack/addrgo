package main

import (
	"github.com/dennstack/addrgo/osm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	go osm.ImportOSMData()

	StartApiServer()
}
