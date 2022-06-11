package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var database *sql.DB

func InitDatabase() {
	var err error
	log.Print("Attempting to connect to database...")
	database, err = sql.Open("postgres", os.Getenv("CONNECTION_STRING"))
	if err != nil {
		log.Fatal("Couldn't connect to database", err)
	}
	log.Print("Database connection successful!")
}

func GetDatabase() *sql.DB {
	return database
}
