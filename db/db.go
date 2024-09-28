package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var (
	once     sync.Once
	instance *sql.DB
	dbErr    error
)

// GetDB returns the singleton instance of the database connection
func GetDB() (*sql.DB, error) {
	once.Do(func() {

		// Get database connection parameters from environment variables
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")

		// Build the DSN (Data Source Name)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error opening database connection: %v", err)
			dbErr = err
			return
		}

		// Set any database configurations (optional)
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * 60 * 60) // 5 hours

		// Ping the database to verify connection
		if err := db.Ping(); err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
			dbErr = err
			return
		}

		instance = db
	})

	return instance, dbErr
}
