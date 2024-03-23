package main

import (
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var chatID int64

const (
	TOKEN = ""
)

var answers = []string{
	"Бесспорно",
	"Предрешено",
	"Без сомнений",
	"Определённо да",
	"Можешь быть уверен в этом",
	"Мне кажется — да",
	"Вероятнее всего",
	"Хорошие перспективы",
	"Знаки говорят — да",
	"Да",
	"Пока не ясно, попробуй снова",
	"Спроси позже",
	"Лучше не рассказывать",
	"Сейчас нельзя предсказать",
	"Сконцентрируйся и спроси опять",
	"Даже не думай",
	"Мой ответ — нет",
	"По моим данным — нет",
	"Перспективы не очень хорошие",
	"Очень сомнительно",
	"Судьба неизбежна",
	"Звезды говорят — да",
	"Сейчас всё складывается удачно",
	"Будущее туманно, спроси еще",
	"Слишком рано, чтобы сказать",
	"Да, но это потребует времени",
	"Возможно, но будь готов к сюрпризу",
	"Сила говорит — да",
	"Не на своем месте сейчас. Переспроси",
	"Ты уже знаешь ответ",
	"Сосредоточься на другом",
	"Ищи знаки вокруг",
	"Да, но не ожидай скоро",
	"Необходимо изменение пути",
	"Сначала исправь ошибки",
	"Будь открыт к новому",
	"Нет, но это изменится",
	"Терпение принесет ответ",
	"Время покажет",
	"Отпусти и получишь ответ",
	"Духи говорит — да",
}

func connectWithTelegram() {
	var err error
	if bot, err = tgbotapi.NewBotAPI(TOKEN); err != nil {
		panic("cannot connect to telegram")
	}
}

func sendMessage(msg string) {
	msgConfig := tgbotapi.NewMessage(chatID, msg)
	bot.Send(msgConfig)
}

func isMessageForTheBot(update *tgbotapi.Update) bool {
	var msg = strings.ToLower(update.Message.Text)
	if update.Message == nil || msg == "" {
		return false
	}
	if strings.Contains(msg, "что будет") {
		return true
	}
	return false
}

func sendAnswer(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(chatID, generateAnswer())
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func generateAnswer() string {
	indx := rand.Intn(len(answers))
	return answers[indx]
}

func main() {
	connectWithTelegram()
	updateConfig := tgbotapi.NewUpdate(0)
	for update := range bot.GetUpdatesChan(updateConfig) {
		chatID = update.Message.Chat.ID
		if update.Message != nil && update.Message.Text == "/start" {
			sendMessage("задай вопрос начав со слов: что будет")
		}
		if isMessageForTheBot(&update) {
			sendAnswer(&update)
		}
	}
}
