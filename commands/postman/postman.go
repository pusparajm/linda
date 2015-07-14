package postman

import (
	"errors"
	"fmt"

	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/kpashka/linda/config"
)

type Postman struct {
	cache map[string]int
	cfg   config.Command
}

func New(cfg config.Command) *Postman {
	c := new(Postman)
	c.cache = map[string]int{}
	c.cfg = cfg
	return c
}

// Return response
func (c *Postman) Run(params []string) (string, error) {
	feed := rss.New(10, false, nil, nil)

	if err := feed.Fetch(c.cfg.Url, nil); err != nil {
		return "", err
	}

	if len(feed.Channels) == 0 {
		return "", errors.New("No items found in feed")
	}

	var response string
	for _, item := range feed.Channels[0].Items {
		if len(item.Links) == 0 {
			continue
		}

		url := item.Links[0].Href
		if _, ok := c.cache[url]; !ok {
			response = fmt.Sprintf("%s: %s", item.Title, url)
			c.cache[url] = 1
			break
		}
	}

	if response == "" {
		return "", errors.New("You've read all news. Please update feed later.")
	}

	return response, nil
}
