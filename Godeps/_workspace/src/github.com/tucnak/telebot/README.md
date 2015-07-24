# Telebot
>Telebot is a convenient wrapper to Telegram Bots API, written in Golang.

[![GoDoc](https://godoc.org/github.com/tucnak/telebot?status.svg)](https://godoc.org/github.com/tucnak/telebot)

Bots are special Telegram accounts designed to handle messages automatically. Users can interact with bots by sending them command messages in private or group chats. These accounts serve as an interface for code running somewhere on your server.

Telebot offers a convenient wrapper to Bots API, so you shouldn't even care about networking at all. Here is an example "helloworld" bot, written with telebot:

```go
import (
    "time"
    "github.com/tucnak/telebot"
)

func main() {
    bot, err := telebot.NewBot("SECRET TOKEN")
    if err != nil {
        return
    }

    messages := make(chan telebot.Message)
    bot.Listen(messages, 1*time.Second)

    for message := range messages {
        if message.Text == "/hi" {
            bot.SendMessage(message.Chat,
                "Hello, "+message.Sender.FirstName+"!", nil)
        }
    }
}
```

You can also send any kind of resources from file system easily:

```go
boom, err := telebot.NewFile("boom.ogg")
if err != nil {
    return err
}

// Next time you send &boom, telebot won't issue
// an upload, but would re-use existing file.
err = bot.SendAudio(recipient, &boom, nil)
```

Sometimes you might want to send a little bit complicated messages, with some optional parameters:

```go
// Send a selective force reply message.
bot.SendMessage(user, "pong", &telebot.SendOptions{
        ForceReply: telebot.ForceReply{
            Require: true,
            Selective: true,
        },
    },
)
```
