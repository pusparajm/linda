package main

import (
	"errors"
	"regexp"
	"strings"
)

func findMatches(expression, input string) ([]string, error) {
	r, err := regexp.Compile(expression)
	if err != nil {
		return []string{}, err
	}

	matches := r.FindStringSubmatch(input)

	if len(matches) == 0 {
		return []string{}, errors.New("No matches found")
	}

	return matches, nil
}

func containsToken(text string, tokens []string) (string, bool) {
	if len(tokens) == 0 {
		return "", false
	}

	toLower := strings.ToLower(text)
	for _, token := range tokens {
		// Temporary - till I'll learn Armenian
		if strings.Contains(text, token) {
			return token, true
		}

		if strings.Contains(toLower, token) {
			return token, true
		}
	}

	return "", false
}
