package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

type WikiResponse struct {
	Query struct {
		Pages map[string]map[string]interface{} `json:"pages,omitempty"`
	} `json:"query,omitempty"`
}

type WikiCommand struct {
	config CmdConfig
	query  string
}

func NewWikiCommand(cfg CmdConfig) *WikiCommand {
	c := new(WikiCommand)
	c.config = cfg
	return c
}

func (c *WikiCommand) GetName() string {
	return c.config.Name
}

func (c *WikiCommand) Trigger(d *DumbSlut, msg *slack.MessageEvent) bool {
	_, ok := containsToken(msg.Text, c.config.Tokens)
	if !ok {
		return false
	}

	return true
}

func (c *WikiCommand) Execute(d *DumbSlut, msg *slack.MessageEvent) {
	serviceResponse, err := c.getWiki(c.query)
	if err != nil {
		log.Error(err.Error())
		d.TalkTo(err.Error(), msg.UserId)
		return
	}

	message := ""
	for _, page := range serviceResponse.Query.Pages {
		rawUrl := page["fullurl"].(string)
		pageUrl, err := url.QueryUnescape(rawUrl)
		if err != nil {
			log.Error(err.Error())
			d.TalkTo(err.Error(), msg.UserId)
			return
		}

		message = fmt.Sprintf("%s: %s", page["title"].(string), pageUrl)
		break
	}

	d.TalkTo(message, msg.UserId)
}

func (c *WikiCommand) getWiki(query string) (*WikiResponse, error) {
	requestUrl := c.config.Url

	log.Infof("Requesting url: %s", requestUrl)
	response, err := http.Get(requestUrl)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error requesting external API")
		return nil, err
	}

	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error reading response body")
		return nil, err
	}

	wiki := new(WikiResponse)
	err = json.Unmarshal(bytes, wiki)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"body":  string(bytes),
		}).Error("Error decoding response body")
		return nil, err
	}

	return wiki, nil
}
