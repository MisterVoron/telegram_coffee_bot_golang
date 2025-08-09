package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func Init() *sql.DB {
	db, err := sql.Open("sqlite", "coffee.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        user_id INTEGER PRIMARY KEY,
        stamp_count INTEGER NOT NULL DEFAULT 0
    );`)
	if err != nil {
		log.Fatal("Create table failed:", err)
	}

	return db
}
