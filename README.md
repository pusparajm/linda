# Linda

[![Build Status](https://travis-ci.org/kpashka/linda.svg)](https://travis-ci.org/kpashka/linda) [![GoDoc](https://godoc.org/github.com/kpashka/linda?status.svg)](https://godoc.org/github.com/kpashka/linda)

Multi-platform, highly configurable conference bot.

Example usage:

<p align="center">
	<img src="http://i.imgur.com/cDKo8FA.png">
</p>

Navigation:

1. [Features](#features)
1. [Limitations](#limitations)
1. [Installation](#installation)
1. [Configuration](#configuration)
1. [Credits](#credits)

## Features

* Different adapters support:
	* [Slack](https://api.slack.com/bot-users)
	* [Telegram](https://core.telegram.org/bots/api)<sup>beta</sup> 
* Configurable commands:
	* [`bully`](commands/bully) - reacts with predefined phrase to matched text.
	* [`copycat`](commands/copycat) - returns same text, can be powerful in combination with filters.
	* [`help`](commands/help) - prints information about instantiated commands.
	* [`postman`](commands/postman) - grabs latest unread item from RSS/Atom feed.
	* [`proxy`](commands/proxy) - fetches JSON document from URL, returns computed template.
* Built-in filters:
	* `base64` - encodes input text to base64.
	* `md5` - calculates md5 sum of input text.
	* `translit` - transliterates input text.
	* `uppercase` - converts input text to uppercase.
* User-friendly:
	* Configurable greeting and farewell messages.
	* Configurable reaction to user status change.
	* Has option to get configuration file by provided URL.
	* `shy mode` in case of being annoyed by large amount of greetings.

## Limitations

Because of the fact that different chat services have different protocols and available options, some usage limitations are present. The table of differences lies below:

| Feature                               | [Slack](https://api.slack.com/bot-users) | [Telegram](https://core.telegram.org/bots/api) |
| ------------------------------------- | ---------------------------------------- | ---------------------------------------------- |
| [`bully`](commands/bully) command     | :white_check_mark: Supported             | :white_check_mark: Supported                   |
| [`copycat`](commands/copycat) command | :white_check_mark: Supported             | :white_check_mark: Supported                   |
| [`help`](commands/help) command   	| :white_check_mark: Supported             | :white_check_mark: Supported	                |
| [`postman`](commands/postman) command | :white_check_mark: Supported             | :white_check_mark: Supported                   |
| [`proxy`](commands/proxy) command     | :white_check_mark: Supported             | :white_check_mark: Supported                   |
| Greetings & farewells                 | :white_check_mark: Supported             | :x: (TBD)                                      |
| Status change reactions               | :white_check_mark: Supported             | :x: (TBD)                                      |

## Installation

Install [godep](https://github.com/tools/godep) tool.

Build and run:

	$ go get https://github.com/kpashka/linda
	$ cd $GOPATH/src/github.com/kpashka/linda
	$ godep restore
	$ go build && ./linda -c <path_to_your_configuration_file>

## Configuration

* See [linda.example.toml](linda.example.toml) for configuration example.
* Detailed explanation for each option available on [wiki](https://github.com/kpashka/linda/wiki/Configuration).

## Credits

* [BurntSushi/toml](https://github.com/BurntSushi/toml)
* [fiam/gounidecode](https://github.com/fiam/gounidecode)
* [jteeuwen/go-pkg-rss](https://github.com/jteeuwen/go-pkg-rss)
* [nlopes/slack](https://github.com/nlopes/slack)
* [NodePrime/jsonpath](https://github.com/NodePrime/jsonpath)
* [Sirupsen/logrus](https://github.com/Sirupsen/logrus)
* [tools/godep](https://github.com/tools/godep)
* [tucnak/telebot](https://github.com/tucnak/telebot)
