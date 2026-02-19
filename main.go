package main

import (
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if token == "" {
		panic("missing TELEGRAM_BOT_TOKEN environment variable")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic("cannot connect to telegram: " + err.Error())
	}
	botUsername = bot.Self.UserName

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	for update := range bot.GetUpdatesChan(updateConfig) {
		handleUpdate(bot, update)
	}
}
