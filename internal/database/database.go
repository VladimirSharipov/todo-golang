package db

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("postgres", "user=root dbname=todo-golang sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	DB.Close()
}
