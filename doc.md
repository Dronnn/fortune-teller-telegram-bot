# fortune-teller.go guide

this document explains the go basics, keywords, and the exact flow of `fortune-teller.go` for a first-time go reader.

## 1) what the bot does (high level)

- connects to telegram using a bot token from the environment
- listens for new messages using long polling
- if a user sends `/start`, it replies with a short help message
- if a message contains the phrase `что будет` (case-insensitive), it replies with a random answer from the list

## 2) required config

- env var: `TELEGRAM_BOT_TOKEN`
  - this is the only required config
  - if it is missing, the program calls `panic` and exits

## 3) how to run (local)

```bash
export TELEGRAM_BOT_TOKEN="your_token_here"
go run fortune-teller.go
```

## 4) dependencies and libraries

**standard library:**

- `math/rand`  
  used to pick a random answer
- `os`  
  used to read the environment variable
- `strings`  
  used to trim and lowercase the incoming text and check for a substring

**third-party:**

- `github.com/go-telegram-bot-api/telegram-bot-api/v5` (alias: `tgbotapi`)  
  provides the telegram bot client and update loop

**go module info:**

- go version: `go 1.22.1`
- module name: `fortune-teller`

## 5) go keywords and syntax (with meanings)

**keywords used in the file:**

- `package`  
  declares the package; `package main` means this file builds an executable
- `import`  
  brings other packages into this file
- `var`  
  declares variables; at package level they live for the whole program
- `func`  
  defines a function
- `if`  
  runs a block only when a condition is true
- `for`  
  the only loop keyword in go
- `range`  
  used with `for` to iterate over slices, maps, or channels
- `return`  
  exits a function (optionally with a value)

**other go keywords (not all are used here):**

- `const`, `type`, `struct`, `interface`
- `switch`, `case`, `default`, `else`
- `break`, `continue`, `goto`, `fallthrough`
- `defer`, `go`, `select`

**not keywords, but commonly confused:**

- `true`, `false`, `nil` (predeclared identifiers)
- `panic`, `recover`, `make`, `new`, `append`, `len`, `cap` (built-in functions)

**important syntax used:**

- `*t`  
  pointer to type `t` (for example `*tgbotapi.BotAPI`)
- `[]t`  
  slice of type `t` (for example `[]string`)
- `:=`  
  short variable declaration inside a function
- `=`  
  assignment to an existing variable
- `.`  
  selector for fields and methods (for example `update.Message`)
- `()`  
  function call with arguments

## 6) telegram-specific types used in the file

- `tgbotapi.BotAPI`  
  the main client used to talk to telegram
- `tgbotapi.Update`  
  a single update from telegram (message, edit, etc.)
- `tgbotapi.NewBotAPI(token)`  
  creates the client
- `tgbotapi.NewUpdate(0)`  
  starts long polling from the latest updates
- `bot.GetUpdatesChan(config)`  
  returns a channel that yields updates as they arrive
- `tgbotapi.NewMessage(chatID, text)`  
  creates a message to send

## 7) step-by-step flow of fortune-teller.go

1. `main()` starts the program.
2. `connectWithTelegram()` reads `TELEGRAM_BOT_TOKEN` and creates the bot client.
3. `NewUpdate(0)` creates a polling config (offset 0 means “from latest”).
4. `GetUpdatesChan` opens long polling and returns a channel of updates.
5. `for update := range channel` loops over incoming updates.
6. the current `chatID` is set from the message.
7. if the message is `/start`, the bot sends a help prompt.
8. if the text contains `что будет`, the bot replies with a random answer.

## 7.1) what is long polling (and alternatives)

**long polling:**

- the client sends a request to telegram and the server keeps the connection open until a new update arrives or a timeout occurs
- once the response arrives, the client immediately opens the next request
- this reduces CPU and network usage compared to very frequent short polls
- in this bot, `GetUpdatesChan` handles the long polling loop for you

**alternatives:**

- **webhooks**  
  telegram sends updates to your HTTPS endpoint. you need a public domain and TLS. lower latency and less outgoing traffic, but requires server setup and a public URL.
- **short polling**  
  client asks for updates on a fixed interval (for example, every 1–2 seconds). simplest to understand but wastes network calls and can be rate-limited.
- **third-party platforms**  
  run the bot on a managed service that handles webhooks and scaling for you (less control, more convenience).

## 8) behavior details and edge cases

- **random answers**  
  uses `math/rand` without a seed, so the sequence is repeatable after each restart.
- **global chat id**  
  `chatID` is global and overwritten on every update. this works for a single chat but can send replies to the wrong chat if multiple chats send messages at the same time.
- **nil message risk**  
  `update.Message` can be `nil` (for example, callback queries). the current code sets `chatID` before checking `update.Message`, which can panic if a non-message update arrives.

## 9) what each function is responsible for

- `connectWithTelegram()`  
  reads the token and initializes the bot client
- `sendMessage(msg)`  
  sends a plain message to the current chat id
- `isMessageForTheBot(update)`  
  checks if message contains `что будет`
- `sendAnswer(update)`  
  replies to the original message with a random answer
- `generateAnswer()`  
  picks a random string from `answers`
- `main()`  
  coordinates everything: connect, poll updates, respond

## 10) quick mental model

think of the bot as a loop that waits for new messages and reacts with two rules:

- rule 1: `/start` -> help text
- rule 2: message contains `что будет` -> random reply

everything else is ignored.
