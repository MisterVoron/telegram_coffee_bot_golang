package bot

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/MisterVoron/telegram_coffee_bot_golang/types"
)

var (
	botApi  string
	adminID string
)

func Start(db *sql.DB) {
	token := os.Getenv("BOT_TOKEN")
	botApi = fmt.Sprintf("https://api.telegram.org/bot%s", token)
	adminID = os.Getenv("ADMIN_ID")
	if token == "" {
		log.Fatal("BOT_TOKEN not set")
	}

	var offset int64

	for {
		updates := getUpdates(offset)
		for _, u := range updates {
			offset = u.UpdateID + 1
			if u.Message != nil {
				handleMessage(db, u.Message)
			}

			if u.CallbackQuery != nil {
				handleCallback(db, u.CallbackQuery)
			}
		}
	}
}

func getUpdates(offset int64) []types.Update {
	url := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=30", botApi, offset)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("GetUpdates failed:", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result types.UpdateResponse
	json.Unmarshal(body, &result)
	return result.Result
}
