package db

import (
	"database/sql"
	"log"
)

func GetStampCount(db *sql.DB, userID int64) int {
	var count int
	err := db.QueryRow("SELECT stamp_count FROM users WHERE user_id = ?", userID).Scan(&count)
	if err == sql.ErrNoRows {
		return 0
	} else if err != nil {
		log.Println("Query error:", err)
		return 0
	}
	return count
}

func IncrementStamp(db *sql.DB, userID int64) int {
	count := GetStampCount(db, userID) + 1

	_, err := db.Exec(`
        INSERT INTO users (user_id, stamp_count)
        VALUES (?, ?)
        ON CONFLICT(user_id) DO UPDATE SET stamp_count = excluded.stamp_count;
    `, userID, count)

	if err != nil {
		log.Println("Update error:", err)
	}
	return count
}

func ResetStamp(db *sql.DB, userID int64) {
	_, err := db.Exec("UPDATE users SET stamp_count = 0 WHERE user_id = ?", userID)
	if err != nil {
		log.Println("Reset error:", err)
	}
}
