package main

import (
  "time"

  "github.com/thoj/go-ircevent"
)

func RejoinCommand(event *irc.Event) {
	go func(event *irc.Event) {
    time.Sleep(time.Second * 3)
		irccon.Join(channel)
	}(event)
}
