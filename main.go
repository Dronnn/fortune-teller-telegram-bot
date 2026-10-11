package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// tokenLogger redacts the bot token from library log output — library errors
// can embed it in URLs (e.g. "<api>/bot<TOKEN>/getUpdates").
type tokenLogger struct {
	token string
}

func (l tokenLogger) Println(v ...interface{}) {
	log.Println(redact(l.token, fmt.Sprint(v...)))
}

func (l tokenLogger) Printf(format string, v ...interface{}) {
	log.Println(redact(l.token, fmt.Sprintf(format, v...)))
}

func redact(token, msg string) string {
	return strings.ReplaceAll(msg, token, "<redacted>")
}

// loadToken prefers the TELEGRAM_BOT_TOKEN env var, then falls back to
// config.env (KEY=VALUE lines) in the working directory.
func loadToken() string {
	if token := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")); token != "" {
		return token
	}
	return tokenFromFile("config.env")
}

func tokenFromFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(key) != "TELEGRAM_BOT_TOKEN" {
			continue
		}
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		return strings.TrimSpace(value)
	}
	return ""
}

func main() {
	token := loadToken()
	if token == "" {
		fmt.Fprintln(os.Stderr, "missing TELEGRAM_BOT_TOKEN: set the env var or put TELEGRAM_BOT_TOKEN=... in config.env")
		os.Exit(1)
	}

	tgbotapi.SetLogger(tokenLogger{token: token})

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(redact(token, "cannot connect to telegram: "+err.Error()))
	}
	botUsername = bot.Self.UserName

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	for update := range bot.GetUpdatesChan(updateConfig) {
		handleUpdate(bot, update)
	}
}
