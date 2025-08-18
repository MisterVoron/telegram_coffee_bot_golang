package bot

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/MisterVoron/telegram_coffee_bot_golang/db"
	"github.com/MisterVoron/telegram_coffee_bot_golang/types"
)

const (
	stampGoal = 6
	validCode = "coffee2025" // Temporary: replace with QR-based code later
)

func handleMessage(database *sql.DB, msg *types.Message) {
	userID := msg.From.ID //userID coincides with chatID
	text := msg.Text

	switch {
	case text == "/start":
		sendMessage(userID, "👋 Welcome to the Coffee Club!\nUse /stamp <code> after each visit.\nCollect 6 stamps to earn a free coffee!")

	case text == "/status":
		count := db.GetStampCount(database, userID)
		sendMessage(userID, fmt.Sprintf("☕ You have %d/%d stamps.", count, stampGoal))

	case strings.HasPrefix(text, "/stamp "):
		code := strings.TrimSpace(text[7:])

		if code != validCode {
			sendMessage(userID, "❌ Invalid code. Try again.")
			return
		}

		text := fmt.Sprintf("☕ Customer @%s requests a stamp", msg.From.Username)
		requestApprovalFromAddmin(text, adminID, userID)
	}
}

func sendMessage(chatID int64, text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	payload := map[string]any{
		"chat_id": chatID,
		"text":    text,
	}

	data, _ := json.Marshal(payload)

	_, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Failed to send message:", err)
	}
}

func requestApprovalFromAddmin(text, adminID string, customerID int64) {
	approve := map[string]any{
		"text":          "✅ Approve",
		"callback_data": fmt.Sprintf("approve:%d", customerID),
	}
	reject := map[string]any{
		"text":          "❌ Reject",
		"callback_data": fmt.Sprintf("reject:%d", customerID),
	}
	keyboard := map[string]any{
		"inline_keyboard": [][]map[string]any{
			{approve, reject},
		},
	}

	kbJSON, _ := json.Marshal(keyboard)
	data := url.Values{}
	data.Set("chat_id", adminID)
	data.Set("text", text)
	data.Set("reply_markup", string(kbJSON))

	http.PostForm(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token), data)
}

func handleCallback(database *sql.DB, callback *types.CallbackQuery) {
	parts := strings.Split(callback.Data, ":")
	if len(parts) != 2 {
		return
	}

	removeInlineKeyboard(callback.Message.Chat.ID, callback.Message.MessageID)

	action := parts[0]
	userID, _ := strconv.Atoi(parts[1])

	if action == "approve" {
		count := db.IncrementStamp(database, int64(userID))
		if count >= stampGoal {
			db.ResetStamp(database, int64(userID))
			sendMessage(int64(userID), "🎉 You earned a FREE coffee! Show this to the barista.")
		} else {
			sendMessage(int64(userID),
				"✅ Stamp added! You now have "+strconv.Itoa(count)+"/"+strconv.Itoa(stampGoal)+" stamps.")
		}
		answerCallback(token, callback.ID, "Stamp approved ✅")
	} else {
		answerCallback(token, callback.ID, "Stamp rejected ❌")
	}
}

func answerCallback(token, callbackID, text string) {
	data := url.Values{}
	data.Set("callback_query_id", callbackID)
	data.Set("text", text)

	http.PostForm(fmt.Sprintf("https://api.telegram.org/bot%s/answerCallbackQuery", token), data)
}

func removeInlineKeyboard(chatID int64, messageID int64) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/editMessageReplyMarkup", token)
	payload := map[string]any{
		"chat_id":      chatID,
		"message_id":   messageID,
		"reply_markup": map[string]any{}, // empty markup removes keyboard
	}
	data, _ := json.Marshal(payload)
	http.Post(url, "application/json", bytes.NewBuffer(data))
}
