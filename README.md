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

The bot token is read from the `TELEGRAM_BOT_TOKEN` env var, or from
`config.env` (gitignored) in the working directory:

```bash
cp config.env.example config.env   # then put your token inside
go run .
```

## Deployment

See `deployment/` for service configs:
- `com.fortune-teller.bot.plist` — launchd (macOS)
- `fortune-teller-telegram-bot.service` — systemd (Linux)

Both run with the repo folder as the working directory and pick the token up
from `config.env` — no token in the service files.

## Project layout

```
main.go            — entrypoint, bot init, update loop
answers.go         — answer categories, weighted random selection, emoji
handlers.go        — trigger phrases, command handlers, response logic
go.mod / go.sum    — module and dependencies
deployment/        — service configs (launchd + systemd)
```
