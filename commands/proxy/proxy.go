package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/NodePrime/jsonpath"
	log "github.com/Sirupsen/logrus"
	"github.com/kpashka/linda/commons"
	"github.com/kpashka/linda/config"
)

// Proxy command - performs HTTP request to JSON API and returns result
type Proxy struct {
	id  string
	cfg config.Command
}

// Creates new instance of proxy command
func New(id string, cfg config.Command) *Proxy {
	c := new(Proxy)
	c.id = id
	c.cfg = cfg
	return c
}

// Runs service command
func (c *Proxy) Run(user *commons.User, params []string) (string, error) {
	// Build request URL
	requestUrl, err := c.buildRequestUrl(params)
	if err != nil {
		return "", err
	}

	// Perform HTTP request
	jsonResponse, err := c.doRequest(requestUrl)
	if err != nil {
		return "", err
	}

	// Format response
	response, err := c.formatResponse(user, jsonResponse)
	if err != nil {
		return "", err
	}

	return response, nil
}

// Builds request url
func (c *Proxy) buildRequestUrl(params []string) (string, error) {
	requestUrl := c.cfg.Url
	for id, param := range params {
		placeholder := "$" + strconv.Itoa(id)
		requestUrl = strings.Replace(requestUrl, placeholder, param, -1)
	}

	u, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
	}

	u.RawQuery = u.Query().Encode()
	return u.String(), nil
}

// Performs http request
func (c *Proxy) doRequest(url string) ([]byte, error) {
	log.WithField("command", c.id).Infof("Requesting URL: %s", url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Response status is not OK")
	}

	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, err
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, err
	}

	body, err = json.Marshal(jsonResponse)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Formats response
func (c *Proxy) formatResponse(user *commons.User, response []byte) (string, error) {
	t, err := template.New("Service Template").Funcs(FuncMap).Parse(c.cfg.Response)
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{}
	for name, selector := range c.cfg.Params {
		// Parse JsonPath expression
		paths, err := jsonpath.ParsePaths(selector)
		if err != nil {
			return "", err
		}

		// Eval JsonPath expression
		eval, err := jsonpath.EvalPathsInBytes(response, paths)
		if err != nil {
			return "", err
		}

		// Find matches
		var result *jsonpath.Result
		for {
			if try, ok := eval.Next(); ok {
				result = try
			} else {
				break
			}
		}

		if result != nil {
			param := string(result.Value)

			// Workaround for lexer strings
			if result.Type == jsonpath.JsonString {
				param = param[1 : len(param)-1]
			}

			params[name] = param
		}
	}

	// Override params
	params["username"] = user.Username
	params["nickname"] = user.Nickname

	// Execute template to string
	doc := bytes.Buffer{}
	err = t.Execute(&doc, params)
	if err != nil {
		return "", err
	}

	return doc.String(), nil
}
