package main

import (
	"log"

	"github.com/MisterVoron/telegram_coffee_bot_golang/bot"
	"github.com/MisterVoron/telegram_coffee_bot_golang/db"
)

func main() {
	database := db.Init()
	defer database.Close() // Close DB on shutdown

	log.Println("Starting Coffee Bot...")
	bot.Start(database)
}
