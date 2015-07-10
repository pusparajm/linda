package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/kpashka/dumbslut/armenia"
	"github.com/nlopes/slack"
)

const (
	MinQueryLength = 2
)

type WeatherResponse struct {
	CityName string `json:"name"`

	Message *string `json:"message,omitempty"`

	Main struct {
		Temperature float64 `json:"temp"`
	} `json:"main"`

	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

type WeatherCommand struct {
	config CmdConfig
	query  string
}

func NewWeatherCommand(cfg CmdConfig) *WeatherCommand {
	c := new(WeatherCommand)
	c.config = cfg
	return c
}

func (c *WeatherCommand) GetName() string {
	return c.config.Name
}

func (c *WeatherCommand) Trigger(d *DumbSlut, msg *slack.MessageEvent) bool {
	currentToken, ok := containsToken(msg.Text, c.config.Tokens)
	if !ok {
		return false
	}

	query := c.normalizeQuery(msg.Text, currentToken)
	if len(query) < MinQueryLength {
		return false
	}

	c.query = query
	return true
}

func (c *WeatherCommand) Execute(d *DumbSlut, msg *slack.MessageEvent) {
	weatherResponse, err := c.getWeather(c.query)
	if err != nil {
		log.Error(err.Error())
		d.Talk(err.Error())
		return
	}

	response, err := c.formatResponse(weatherResponse)
	if err != nil {
		log.Error(err.Error())
		d.Talk(err.Error())
		return
	}

	d.Talk(response)
}

func (c *WeatherCommand) buildRequestUrl() string {
	return c.config.Url + c.query
}

func (c *WeatherCommand) normalizeQuery(text, token string) string {
	log.Infof("Normalizing query by token %s", token)

	if token == "Եղանակը" { // Armenian trigger
		query := strings.Replace(text, token, "", 1)
		queryWords := strings.Split(strings.TrimSpace(query), " ")
		city := armenia.Translit(queryWords[0], armenia.WesternLanguage)
		return city
	}

	query := strings.Replace(strings.ToLower(text), token, "", 1)
	queryWords := strings.Split(strings.TrimSpace(query), " ")
	return queryWords[0]
}

func (c *WeatherCommand) getWeather(query string) (*WeatherResponse, error) {
	requestUrl := c.buildRequestUrl()

	log.Infof("Requesting url: %s", requestUrl)
	response, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	weather := new(WeatherResponse)
	err = json.Unmarshal(bytes, weather)
	if err != nil {
		return nil, err
	}

	if weather.Message != nil {
		return nil, errors.New(fmt.Sprintf("Can't find weather for %s", query))
	}

	return weather, nil
}

func (c *WeatherCommand) formatResponse(response *WeatherResponse) (string, error) {
	t, err := template.New("Weather Template").Parse(c.config.Response)
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{
		"city":        response.CityName,
		"temperature": fmt.Sprintf("%0.0f", response.Main.Temperature),
		"description": response.Weather[0].Description,
	}

	doc := bytes.Buffer{}
	err = t.Execute(&doc, params)
	if err != nil {
		return "", err
	}

	return doc.String(), nil
}
