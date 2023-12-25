package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RSO-project-Prepih/get-photo-info/prometheus"
	_ "github.com/lib/pq"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var DB *sql.DB

func NewDBConnection() *sql.DB {

	pgURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=verify-full",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var db *sql.DB
	var err error

	maxRetries := 5

	delay := 500 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		prometheus.DBConnectionAttempts.Inc()
		db, err = sql.Open("postgres", pgURL)
		if err == nil {
			prometheus.DBConnectionAttempts.Inc()
			err = db.Ping()
			if err == nil {
				log.Println("Connected to database")
				return db
			}
		}

		log.Printf("Failed to connect to database: %v. Retrying in %v...\n", err, delay)
		time.Sleep(delay)

		// Increase the delay for the next retry using exponential backoff
		delay *= 2
	}

	logFatal(err)
	return nil
}

func TestConnection() {
	rows, err := DB.Query("SELECT name FROM users")
	logFatal(err)
	defer rows.Close()

	fmt.Println("Usernames from the 'users' table:")
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		logFatal(err)
		fmt.Println(username)
	}

	err = rows.Err()
	logFatal(err)
}
