package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type (
	Database struct {
		Connection *sql.DB
	}
)

var (
	instance *Database
	once     sync.Once
)

func GetDatabaseInstance() *Database {
	once.Do(func() {
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASS")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		dbname := os.Getenv("DB_NAME")

		connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)
		db, err := sql.Open("mysql", connStr)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(10)
		db.SetConnMaxLifetime(5 * time.Minute)

		err = db.Ping()
		if err != nil {
			log.Fatal("Failed to ping database:", err)
		}
		log.Println("Successfully connected to the database")

		instance = &Database{Connection: db}
	})
	return instance
}
