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

* Different adapters support:
	* [`slack`](adapters/slack)
	* [`telegram`](adapters/telegram)<sup>beta</sup> 
* Configurable commands:
	* [`artist`](command/artist) - draws symbolic ASCII art from input word.
	* [`bully`](command/bully) - reacts with pre-defined phrase to matched text.
	* [`copycat`](command/copycat) - returns same text, can be powerful in combination with filters.
	* [`postman`](command/postman) - grabs latest unread item from RSS/Atom feed.
	* [`proxy`](command/proxy) - fetches JSON document from URL, returns computed template.
* Built-in commands:
	* [`help`](command/help) - prints information about instanced commands.
* Built-in filters:
	* `base64` - encodes input text to base64.
	* `md5` - calculates md5 sum of input text.
	* `translit` - transliterates input text.
	* `uppercase` - converts input text to uppercase.
* User-friendly:
	* Configurable greeting and farewell messages.
	* Configurable reaction to user status change.
	* "Shy mode" in case of being annoyed by chatterbox.

## Limitations

Because of the fact that different chat services have different protocols and available options, some usage limitations are present. The table of differences lies below:

| Feature                              | [`slack`](adapters/slack)    | [`telegram`](adapters/telegram)    |
| ------------------------------------ | ---------------------------- | ---------------------------------- |
| [`artist`](command/artist) command   | :white_check_mark: Supported | :x: (no Markdown support)          |
| [`bully`](command/bully) command     | :white_check_mark: Supported | :white_check_mark: Supported       |
| [`copycat`](command/copycat) command | :white_check_mark: Supported | :white_check_mark: Supported       |
| [`help`](command/help) command   	   | :white_check_mark: Supported | :white_check_mark: Supported	   |
| [`postman`](command/postman) command | :white_check_mark: Supported | :white_check_mark: Supported       |
| [`proxy`](command/proxy) command     | :white_check_mark: Supported | :white_check_mark: Supported       |
| Greetings & farewells                | :white_check_mark: Supported | :x: (TBD)                          |
| Status change reactions              | :white_check_mark: Supported | :x: (TBD)                          |

## Installation

Build and run:

	$ mkdir -p $GOPATH/src/github.com/kpashka/linda
	$ cd $GOPATH/src/github.com/kpashka/linda
	$ git clone https://github.com/kpashka/linda .
	$ go get -v ./...
	$ go build && ./linda -c config.toml

## Configuration

* See [config.example.toml](config.example.toml) for configuration example.

## Dependencies

* [`BurntSushi/toml`](https://github.com/BurntSushi/toml)
* [`fiam/gounidecode`](https://github.com/fiam/gounidecode)
* [`jteeuwen/go-pkg-rss`](https://github.com/jteeuwen/go-pkg-rss)
* [`nlopes/slack`](https://github.com/nlopes/slack)
* [`NodePrime/jsonpath`](https://github.com/NodePrime/jsonpath)
* [`Sirupsen/logrus`](https://github.com/Sirupsen/logrus)
* [`tucnak/telebot`](https://github.com/tucnak/telebot)