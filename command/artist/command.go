package artist

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/kpashka/dumbslut/config"
)

const (
	MinLen = 3
	MaxLen = 10
)

type Artist struct {
	cfg          config.Command
	currentToken string
}

func New(cfg config.Command) *Artist {
	c := new(Artist)
	c.cfg = cfg
	return c
}

// Return response
func (c *Artist) Run(params []string) (string, error) {
	if len(params) < 1 {
		return "", errors.New("Matching error")
	}

	inputWord := params[1]
	runeCount := utf8.RuneCountInString(inputWord)
	if runeCount < MinLen || runeCount > MaxLen {
		return "", errors.New(fmt.Sprintf("Input string size should be in range [%d, %d]", MinLen, MaxLen))
	}

	swastika := NewSwastika(inputWord)
	markdown := fmt.Sprintf("\n```%s```", swastika)
	return markdown, nil
}
