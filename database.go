package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	Db *sql.DB
)

func ConnectDB(passwd string) *sql.DB {
	hostname := getenv("DBHOST", "192.168.1.249")
	username := getenv("DBUSER", "thatmaidguy")
	password := passwd
	database := getenv("DBNAME", "hjam")
	sslMode := "disable"

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", username, password, hostname, database, sslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println("db.Ping() показал:")
		log.Fatal(err)
	}

	return db
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
