package main

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var botUsername string

var triggerPhrases = []string{
	"что будет", "будет ли", "стоит ли",
	"предскажи", "скажи", "подскажи",
	"а что если", "что если",
	"суждено ли", "ждёт ли меня", "ждет ли меня",
	"получится ли", "смогу ли", "удастся ли",
	"правда ли", "верно ли",
	"есть ли шанс", "есть ли смысл",
	"может ли", "можно ли", "стоит мне",
	"выйдет ли", "сбудется ли",
	"повезёт ли", "повезет ли",
	"случится ли", "произойдёт ли", "произойдет ли",
	"ожидать ли", "надеяться ли",
	"будет хорошо", "всё будет",
	"магический шар", "шар предсказаний",
	"о великий", "о мудрый",
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	msg := update.Message
	chatID := msg.Chat.ID

	if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			handleStart(bot, chatID)
		case "help":
			handleHelp(bot, chatID)
		}
		return
	}

	if !shouldRespond(msg) {
		return
	}

	reply := tgbotapi.NewMessage(chatID, generateAnswer())
	reply.ReplyToMessageID = msg.MessageID
	bot.Send(reply)
}

func handleStart(bot *tgbotapi.BotAPI, chatID int64) {
	text := "🔮 Я — магический шар предсказаний!\n\nЗадай мне вопрос, и я отвечу.\nНапиши /help чтобы узнать подробнее."
	bot.Send(tgbotapi.NewMessage(chatID, text))
}

func handleHelp(bot *tgbotapi.BotAPI, chatID int64) {
	text := `🔮 *Как спросить предсказание:*

Начни вопрос с одной из фраз:
• "Что будет, если..."
• "Будет ли..."
• "Стоит ли мне..."
• "Предскажи..."
• "Суждено ли..."
• "Повезёт ли..."
• "О великий шар..."

В групповом чате — упомяни @` + botUsername + ` или ответь на моё сообщение.`
	helpMsg := tgbotapi.NewMessage(chatID, text)
	helpMsg.ParseMode = "Markdown"
	bot.Send(helpMsg)
}

func shouldRespond(msg *tgbotapi.Message) bool {
	text := strings.ToLower(msg.Text)
	if text == "" {
		return false
	}

	hasTrigger := containsTrigger(text)
	lowerUsername := strings.ToLower(botUsername)
	mentioned := strings.Contains(text, "@"+lowerUsername) || strings.Contains(text, lowerUsername)
	repliedToBot := msg.ReplyToMessage != nil && msg.ReplyToMessage.From != nil &&
		msg.ReplyToMessage.From.UserName == botUsername

	// Private chat: trigger phrase or bot name
	if msg.Chat.Type == "private" {
		return hasTrigger || mentioned
	}

	// Group chat: @mention or reply to bot is enough; trigger alone is not
	if mentioned || repliedToBot {
		return true
	}
	return false
}

func containsTrigger(text string) bool {
	for _, trigger := range triggerPhrases {
		if strings.Contains(text, trigger) {
			return true
		}
	}
	return false
}
