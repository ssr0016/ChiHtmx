package database

import (
	"database/sql"
	"log"
)

var DBConn *sql.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=secret dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	DBConn = db
	log.Println("Connected to DB")
}
