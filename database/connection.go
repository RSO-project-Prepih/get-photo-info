package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var DB *sql.DB

func init() {
	// Load .env file for environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	pgURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=verify-full",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Open connection to database
	var err error
	DB, err = sql.Open("postgres", pgURL)
	if err = DB.Ping(); err != nil {
		logFatal(err)
	}

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
