# Dumb Slut

[![GoDoc](https://godoc.org/github.com/kpashka/dumbslut?status.svg)](https://godoc.org/github.com/kpashka/dumbslut)

:person_with_pouting_face: Little princess, programmed to serve in mens-only conference rooms.

## Features

* Different backends support:
	* [Slack](https://api.slack.com/bot-users)
	* [Telegram](https://core.telegram.org/bots)<sup>beta</sup> 
* Configurable commands:
	* `Artist` - draws symbolic ASCII art from input word.
	* `Bully` - reacts with pre-defined phrase to matched text.
	* `Postman` - grabs latest unread item from RSS/Atom feed.
	* `Proxy` - fetches JSON document from URL, returns template with populated values from [JSONPath](https://github.com/NodePrime/jsonpath#path-syntax) expressions.
	* `Snitch` - returns an info about of configured command instances.
* User-friendly:
	* Configurable greeting and farewell messages.
	* Configurable reaction to user status change.
	* "Shy" mode in case of being annoyed by chatterbox servant.
	* Live configuration reload from URL - share access to the ear with your mates (upcoming).

## Limitations

Because of the fact that backend services were created by different people, there are some usage limitations. The table of differences lies below:

| Feature                 | [Slack](https://api.slack.com/bot-users) | [Telegram](https://core.telegram.org/bots) |
| ----------------------- | ---------------------------------------- | ------------------------------------------ |
| Greetings               | :white_check_mark: Supported             | :x: (not sure, if possible)                |
| Farewells               | :white_check_mark: Supported             | :x: (probably possible)                    |
| Status change reactions | :white_check_mark: Supported             | :x: (looks like impossible)                |
| `Artist` command        | :white_check_mark: Supported             | :x: (no Markdown support)                  |
| `Bully` command         | :white_check_mark: Supported             | :white_check_mark: Supported               |
| `Proxy` command         | :white_check_mark: Supported             | :white_check_mark: Supported               |
| `Postman` command       | :white_check_mark: Supported             | :white_check_mark: Supported               |
| `Snitch` command        | :white_check_mark: Supported             | :white_check_mark: Limited support         |

## Usage

Grab dependencies:

	$ go get github.com/jteeuwen/go-pkg-rss
	$ go get github.com/nlopes/slack
	$ go get github.com/Sirupsen/logrus
	$ go get github.com/tucnak/telebot
	$ go get golang.org/x/exp

Build and run:

	$ go get github.com/kpashka/dumbslut
	$ cd $GOPATH/src/github.com/kpashka/dumbslut
	$ go build && ./dumpslut -c config.json

## Configuration

See [example configuration](config.example.json) for details.