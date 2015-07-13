# Linda

[![GoDoc](https://godoc.org/github.com/kpashka/linda?status.svg)](https://godoc.org/github.com/kpashka/linda)

:princess: Little princess, programmed to serve in mens-only online conference rooms.

Navigation:

1. [Features](#features)
1. [Limitations](#limitations)
1. [Installation](#installation)
1. [Configuration](#configuration)

## Features

* Different backends support:
	* [`Slack`](backend/slack)
	* [`Telegram`](backend/telegram)<sup>beta</sup> 
* Configurable commands:
	* [`Artist`](command/artist) - draws symbolic ASCII art from input word.
	* [`Bully`](command/bully) - reacts with pre-defined phrase to matched text.
	* [`Postman`](command/postman) - grabs latest unread item from RSS/Atom feed.
	* [`Proxy`](command/proxy) - fetches JSON document from URL, returns computed template.
	* [`Snitch`](command/snitch) - returns an info about of configured command instances.
* User-friendly:
	* Configurable greeting and farewell messages.
	* Configurable reaction to user status change.
	* "Shy" mode in case of being annoyed by chatterbox servant.
	* Live configuration reload from URL - share access to the ear with your mates (upcoming).

## Limitations

Because of the fact that backend services were created by different people, there are some usage limitations. The table of differences lies below:

| Feature                              | [Slack](backend/slack)       | [Telegram](backend/telegram)       |
| ------------------------------------ | ---------------------------- | ---------------------------------- |
| [`Artist`](command/artist) command   | :white_check_mark: Supported | :x: (no Markdown support)          |
| [`Bully`](command/bully) command     | :white_check_mark: Supported | :white_check_mark: Supported       |
| [`Postman`](command/postman) command | :white_check_mark: Supported | :white_check_mark: Supported       |
| [`Proxy`](command/proxy) command     | :white_check_mark: Supported | :white_check_mark: Supported       |
| [`Snitch`](command/snitch) command   | :white_check_mark: Supported | :white_check_mark: Limited support |
| Greetings & farewells                | :white_check_mark: Supported | :x: (TBD)                          |
| Status change reactions              | :white_check_mark: Supported | :x: (TBD)                          |

## Installation

Grab dependencies:

	$ go get github.com/jteeuwen/go-pkg-rss
	$ go get github.com/nlopes/slack
	$ go get github.com/Sirupsen/logrus
	$ go get github.com/tucnak/telebot
	$ go get golang.org/x/exp

Build and run:

	$ go get github.com/kpashka/linda
	$ cd $GOPATH/src/github.com/kpashka/linda
	$ go build && ./linda -c config.json

## Configuration

* See [config.example.json](config.example.json) for configuration example.
* See [Configuration](https://github.com/kpashka/linda/wiki/Configuration) page for detailed information.