# fortune-teller-telegram-bot

Telegram "fortune-teller" bot (Magic 8-Ball style) that replies in Russian.

## How it works

- `/start` — greeting
- `/help` — list of trigger phrases
- 36 trigger phrases: "что будет", "будет ли", "предскажи", "о великий шар", etc.
- 75 answers in 3 categories with Classic Magic 8-Ball probabilities (50% positive, 25% uncertain, 25% negative)
- Works in private chats and group chats (via @mention or reply)

## Stack

- Go 1.22+
- [telegram-bot-api v5](https://github.com/go-telegram-bot-api/telegram-bot-api)

## Setup

```bash
export TELEGRAM_BOT_TOKEN="your_token_here"
go run .
```

## Deployment

See `deployment/` for service configs:
- `com.fortune-teller.bot.plist` — launchd (macOS)
- `fortune-teller-telegram-bot.service` — systemd (Linux)

Replace `YOUR_TOKEN_HERE` with your bot token.

## Project layout

```
main.go            — entrypoint, bot init, update loop
answers.go         — answer categories, weighted random selection, emoji
handlers.go        — trigger phrases, command handlers, response logic
go.mod / go.sum    — module and dependencies
deployment/        — service configs (launchd + systemd)
```
