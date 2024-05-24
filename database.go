package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/tanimutomo/sqlfile"
)

func ConnectDB() *sql.DB {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")
	sslMode := "disable"

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", username, password, "db", database, sslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println("db.Ping() показал:")
		log.Fatal(err)
	}

	if os.Getenv("FIRST_RUN") == "1" {
		s := sqlfile.New()

		// Load input file and store queries written in the file
		err = s.File("sql/hjam.sql")
		if err != nil {
			log.Fatal("Невозможно получить файл")
		}
		_, err = s.Exec(db)
		if err != nil {
			log.Fatal("Невозможно мигрировать базу данных")
		}
	}

	return db
}
