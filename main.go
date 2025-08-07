package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	apiURL    = "https://api.telegram.org/bot"
	stampGoal = 6
	validCode = "coffee2025" // Temporary, to be replaced with QR-code logic
)

var (
	botToken   string
	userStamps = make(map[int64]int)
	mu         sync.Mutex
)

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type SendMessagePayload struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func main() {
	botToken = os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Please set BOT_TOKEN environment variable")
	}

	offset := 0

	for {
		updates, err := getUpdates(offset)
		if err != nil {
			log.Println("Error getting updates:", err)
			continue
		}

		for _, update := range updates {
			offset = update.UpdateID + 1
			handleMessage(update.Message)
		}
	}
}

func getUpdates(offset int) ([]Update, error) {
	url := fmt.Sprintf("%s%s/getUpdates?offset=%d&timeout=30", apiURL, botToken, offset)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updateResp UpdateResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &updateResp)
	if err != nil {
		return nil, err
	}

	return updateResp.Result, nil
}

func sendMessage(chatID int64, text string) {
	url := fmt.Sprintf("%s%s/sendMessage", apiURL, botToken)

	payload := SendMessagePayload{
		ChatID: chatID,
		Text:   text,
	}

	data, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Failed to send message:", err)
		return
	}
	defer resp.Body.Close()
}

func handleMessage(msg Message) {
	text := msg.Text
	chatID := msg.Chat.ID

	switch {
	case text == "/start":
		sendMessage(chatID, "👋 Welcome to the Coffee Club!\nUse /stamp <code> to log your visit.\nCollect 6 to earn a free coffee!")
	case text == "/status":
		mu.Lock()
		count := userStamps[chatID]
		mu.Unlock()

		sendMessage(chatID, fmt.Sprintf("☕ You have %d/%d stamps.", count, stampGoal))
	case strings.HasPrefix(text, "/stamp "):
		code := strings.TrimSpace(text[7:])

		if code != validCode {
			sendMessage(chatID, "❌ Invalid code. Try again.")
			return
		}

		mu.Lock()
		userStamps[chatID]++
		count := userStamps[chatID]
		if count >= stampGoal {
			userStamps[chatID] = 0
			mu.Unlock()
			sendMessage(chatID, "🎉 You earned a FREE coffee! Show this to the barista.")
		} else {
			mu.Unlock()
			sendMessage(chatID, fmt.Sprintf("✅ Stamp added! You now have %d/%d stamps.", count, stampGoal))
		}
	}
}
