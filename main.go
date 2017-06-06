package main

import (
	"crypto/tls"
	"fmt"
	"github.com/exponent-io/jsonpath"
	"github.com/thoj/go-ircevent"
	"os"
	"strings"
	"sync"
)

const channel = "#rtk"
const serverssl = "irc.rizon.net:9999"

var dictionaries = [...]string{
	"daijirin",
	"daijisen",
	"kotowaza",
	"meikyou",
}

func main() {
	// Init
	nickname := "kareebot"
	irccon := irc.IRC(nickname, "coco")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Event handlers
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("PRIVMSG", jjLookup)

	err := irccon.Connect(serverssl)

	// Houston, we've got a problem
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}

	irccon.Loop()
}

func jjLookup(event *irc.Event) {
	go func(event *irc.Event) {
		// Split message into tokens
		tokens := strings.Fields(event.Message())

		if tokens[0] == "!jj" && len(tokens) == 2 {
			failed := 0
			var wg = sync.WaitGroup{}
			wg.Add(len(dictionaries))
			failedLock := sync.RWMutex{}

			for _, dictionary := range dictionaries {
				go func(dict string, wg *sync.WaitGroup) {
					result, err := FindInDictionary(tokens[1], dict)

					// Unexpected errors
					if err != nil {
						fmt.Printf("ERROR: %s\n", err)

						failedLock.WLock()
						failed++
						failedLock.WUnlock()

						wg.Done()
						return
					}

					// Not found
					if result == "" {
						failedLock.WLock()
						failed++
						failedLock.WUnlock()

						wg.Done()
						return
					}

					// Post results
					irccon.Privmsg("#rtk", fmt.Sprintf("%s: %s", dict, result))
					wg.Done()
				}(dictionary, &wg)
			}

			wg.Wait()

			// Post in case nothing is found
			failedLock.RLock()
			if failed == 4 {
				irccon.Privmsg("#rtk", "Nope.")
			}
			failedLock.RUnlock()
		}
	}(event)
}

// FindInDictionary searches for a given word inside a JSON dictionary file
// with structure `"word": "definition"`
func FindInDictionary(word string, dictionary string) (string, error) {
	var result string
	file, e := os.Open(fmt.Sprintf("./dicts/%s.json", dictionary))

	if e != nil {
		return "", fmt.Errorf("file error: %v\n", e)
	}

	w := jsonpath.NewDecoder(file)
	w.SeekTo(word)
	w.Decode(&result)

	return result, nil
}
