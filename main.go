package main

import (
	"crypto/tls"
	"fmt"

	"github.com/thoj/go-ircevent"
)

const channel = "#rtk"
const serverssl = "irc.rizon.net:9999"

var irccon *irc.Connection

func main() {
	// Init
	nickname := "kareebot"
	irccon = irc.IRC(nickname, "coco")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Event handlers
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("PRIVMSG", JJLookupCommand)
	irccon.AddCallback("PRIVMSG", TobuCommand)

	err := irccon.Connect(serverssl)

	// Houston, we've got a problem
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}

	irccon.Loop()
}
