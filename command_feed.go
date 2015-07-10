package main

import (
	"fmt"
	"math/rand"

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

func (c *FeedCommand) GetType() string {
	return CommandTypeFeed
}

func (c *FeedCommand) Trigger(d *DumbSlut, msg *slack.MessageEvent) bool {
	_, ok := containsToken(msg.Text, c.config.Tokens)
	return ok
}

func (c *FeedCommand) Execute(d *DumbSlut, msg *slack.MessageEvent) {
	feed := rss.New(10, false, nil, nil)

	if err := feed.Fetch(c.config.ApiUrl, nil); err != nil {
		log.Error(err.Error())
		d.Talk(err.Error())
	}

	response := ""
	if len(feed.Channels) > 0 {
		itemId := rand.Intn(len(feed.Channels[0].Items))
		item := feed.Channels[0].Items[itemId]
		response = fmt.Sprintf("%s: %s", item.Title, item.Links[0].Href)
	}

	d.Talk(response)
}
