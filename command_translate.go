package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

const (
	TranslateMinQuery = 2
)

type TranslateResponse struct {
	Matches []struct {
		Translation string `json:"translation,omitempty"`
		Reference   string `json:"reference,omitempty"`
	} `json:"matches,omitempty"`
}

type TranslateCommand struct {
	config CmdConfig
	query  string
}

func NewTranslateCommand(cfg CmdConfig) *TranslateCommand {
	c := new(TranslateCommand)
	c.config = cfg
	return c
}

func (c *TranslateCommand) GetName() string {
	return c.config.Name
}

func (c *TranslateCommand) Trigger(d *DumbSlut, msg *slack.MessageEvent) bool {
	currentToken, ok := containsToken(msg.Text, c.config.Tokens)
	if !ok {
		return false
	}

	query := c.normalizeQuery(msg.Text, currentToken)
	if len(query) < TranslateMinQuery {
		return false
	}

	c.query = query
	return true
}

func (c *TranslateCommand) Execute(d *DumbSlut, msg *slack.MessageEvent) {
	serviceResponse, err := c.getTranslation(c.query)
	if err != nil {
		log.Error(err.Error())
		d.TalkTo(err.Error(), msg.UserId)
		return
	}

	response, err := c.formatResponse(serviceResponse)
	if err != nil {
		log.Error(err.Error())
		d.TalkTo(err.Error(), msg.UserId)
		return
	}

	d.TalkTo(response, msg.UserId)
}

func (c *TranslateCommand) buildRequestUrl() string {
	return c.config.Url + url.QueryEscape(c.query)
}

func (c *TranslateCommand) normalizeQuery(text, token string) string {
	log.Infof("Normalizing query by token %s", token)
	query := strings.Replace(strings.ToLower(text), token, "", 1)
	return strings.TrimSpace(query)
}

func (c *TranslateCommand) getTranslation(query string) (*TranslateResponse, error) {
	requestUrl := c.buildRequestUrl()

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

	translation := new(TranslateResponse)
	err = json.Unmarshal(bytes, translation)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"body":  string(bytes),
		}).Error("Error decoding response body")
		return nil, err
	}

	if len(translation.Matches) == 0 {
		return nil, errors.New(fmt.Sprintf("Can't find translation for %s", query))
	}

	return translation, nil
}

func (c *TranslateCommand) formatResponse(response *TranslateResponse) (string, error) {
	t, err := template.New("Translation Template").Parse(c.config.Response)
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{
		"translation": response.Matches[0].Translation,
		"reference":   response.Matches[0].Reference,
	}

	doc := bytes.Buffer{}
	err = t.Execute(&doc, params)
	if err != nil {
		return "", err
	}

	return doc.String(), nil
}
