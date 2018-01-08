package main

import (
	"fmt"
	"strings"

	"github.com/thoj/go-ircevent"
)

func TobuCommand(event *irc.Event) {
	go func(event *irc.Event) {
		// Split message into tokens
		tokens := strings.Fields(event.Message())

		if len(tokens) == 0 {
			return
		}

		if tokens[0] != "!tobu" {
			return
		}

		name := "Austin"
		if len(tokens) == 2 {
			name = tokens[1]
		}

		irccon.Privmsg("#rtk", fmt.Sprintf("%s tobu off a hashi pls", name))
	}(event)
}
