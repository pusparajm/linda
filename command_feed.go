package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/nlopes/slack"
)

type FeedCommand struct {
	config CmdConfig
}

func NewFeedCommand(cfg CmdConfig) *FeedCommand {
	c := new(FeedCommand)
	c.config = cfg
	return c
}

func (c *FeedCommand) Trigger(d *DumbSlut, msg *slack.MessageEvent) bool {
	_, ok := containsToken(msg.Text, c.config.Tokens)
	return ok
}

func (c *FeedCommand) Respond(d *DumbSlut, msg *slack.MessageEvent) {
	feed := rss.New(10, false, nil, nil)

	if err := feed.Fetch(c.config.ApiUrl, nil); err != nil {
		log.Error(err.Error())
		d.Talk(err.Error())
	}

	response := ""
	for _, channel := range feed.Channels {
		for _, item := range channel.Items {
			response = fmt.Sprintf("%s: %s", item.Title, item.Links[0].Href)
			break
		}
	}

	d.Talk(response)
}
