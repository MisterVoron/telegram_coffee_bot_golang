package db

import (
	"database/sql"
	"log"
)

func GetStampCount(db *sql.DB, chatID int64) int {
	var count int
	err := db.QueryRow("SELECT stamp_count FROM users WHERE chat_id = ?", chatID).Scan(&count)
	if err == sql.ErrNoRows {
		return 0
	} else if err != nil {
		log.Println("Query error:", err)
		return 0
	}
	return count
}

func IncrementStamp(db *sql.DB, chatID int64) int {
	count := GetStampCount(db, chatID) + 1

	_, err := db.Exec(`
        INSERT INTO users (chat_id, stamp_count)
        VALUES (?, ?)
        ON CONFLICT(chat_id) DO UPDATE SET stamp_count = excluded.stamp_count;
    `, chatID, count)

	if err != nil {
		log.Println("Update error:", err)
	}
	return count
}

func ResetStamp(db *sql.DB, chatID int64) {
	_, err := db.Exec("UPDATE users SET stamp_count = 0 WHERE chat_id = ?", chatID)
	if err != nil {
		log.Println("Reset error:", err)
	}
}
