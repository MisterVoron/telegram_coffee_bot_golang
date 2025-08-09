package bot

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/MisterVoron/telegram_coffee_bot_golang/db"
	"github.com/MisterVoron/telegram_coffee_bot_golang/types"
)

const (
	stampGoal = 6
	validCode = "coffee2025" // Temporary: replace with QR-based code later
)

func HandleMessage(database *sql.DB, msg types.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID
	text := msg.Text

	switch {
	case text == "/start":
		SendMessage(chatID, "👋 Welcome to the Coffee Club!\nUse /stamp <code> after each visit.\nCollect 6 stamps to earn a free coffee!")

	case text == "/status":
		count := db.GetStampCount(database, userID)
		SendMessage(chatID, fmt.Sprintf("☕ You have %d/%d stamps.", count, stampGoal))

	case strings.HasPrefix(text, "/stamp "):
		code := strings.TrimSpace(text[7:])

		if code != validCode {
			SendMessage(chatID, "❌ Invalid code. Try again.")
			return
		}

		count := db.IncrementStamp(database, userID)
		if count >= stampGoal {
			db.ResetStamp(database, userID)
			SendMessage(chatID, "🎉 You earned a FREE coffee! Show this to the barista.")
		} else {
			SendMessage(chatID, "✅ Stamp added! You now have "+strconv.Itoa(count)+"/"+strconv.Itoa(stampGoal)+" stamps.")
		}
	}
}
