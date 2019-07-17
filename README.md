# Telelog - Log via Telegram

Sends your log to your Telegram account. Uses the Telegram Bot API.

## Setup

- Create a Bot (and get its token), see https://core.telegram.org/bots
- Start chatting with the bot (and actually send something - needed for the next step)
- To get the `CHAT_ID` send the following request:
`curl https://api.telegram.org/bot<TOKEN>/getUpdates`

## Usage

```go
telelog := telelog.NewTelelog("TOKEN", "CHAT_ID") // os.Getenv("TOKEN")
telelog.logError("Here is some text")
telelog.logAlert("Here is some text")
telelog.logInfo("Here is some text")
telelog.logSuccess("Here is some text")
```

![](output.png)