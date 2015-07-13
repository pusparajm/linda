# Linda

[![Build Status](https://travis-ci.org/kpashka/linda.svg)](https://travis-ci.org/kpashka/linda) [![GoDoc](https://godoc.org/github.com/kpashka/linda?status.svg)](https://godoc.org/github.com/kpashka/linda)

:princess: Little princess, programmed to serve in mens-only online conference rooms.

Navigation:

1. [Features](#features)
1. [Limitations](#limitations)
1. [Installation](#installation)
1. [Configuration](#configuration)
1. [Dependencies](#dependencies)

## Features

* Different backends support:
	* [`Slack`](backend/slack)
	* [`Telegram`](backend/telegram)<sup>beta</sup> 
* Configurable commands:
	* [`Artist`](command/artist) - draws symbolic ASCII art from input word.
	* [`Bully`](command/bully) - reacts with pre-defined phrase to matched text.
	* [`Postman`](command/postman) - grabs latest unread item from RSS/Atom feed.
	* [`Proxy`](command/proxy) - fetches JSON document from URL, returns computed template.
	* [`Snitch`](command/snitch) - prints information about available computed commands.
* User-friendly:
	* Configurable greeting and farewell messages.
	* Configurable reaction to user status change.
	* `shy` mode in case of being annoyed by chatterbox servant.
	* Live configuration reload from URL - share access to the ear with your mates (TBD).

## Limitations

Because of the fact that backend services have different protocols and available options, some usage limitations are present. The table of differences lies below:

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

Build and run:

	$ mkdir -p $GOPATH/src/github.com/kpashka/linda
	$ cd $GOPATH/src/github.com/kpashka/linda
	$ git clone https://github.com/kpashka/linda .
	$ go get -v ./...
	$ go build && ./linda -c config.json

## Configuration

* See [config.example.json](config.example.json) for configuration example.
* See [Configuration](https://github.com/kpashka/linda/wiki/Configuration) page for detailed information.

## Dependencies

* [`jteeuwen/go-pkg-rss`](github.com/jteeuwen/go-pkg-rss)
* [`nlopes/slack`](github.com/nlopes/slack)
* [`NodePrime/jsonpath`](github.com/NodePrime/jsonpath)
* [`Sirupsen/logrus`](github.com/Sirupsen/logrus)
* [`tucnak/telebot`](github.com/tucnak/telebot)