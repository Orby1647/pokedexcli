package main

import (
	"strings"
)

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	if text == "" {
		return []string{}
	}

	words := strings.Fields(text)
	return words
}
