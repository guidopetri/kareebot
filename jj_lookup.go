package main

import (
  "fmt"
  "os"
  "strings"
  "sync"

  "github.com/exponent-io/jsonpath"
  "github.com/thoj/go-ircevent"
)

var dictionaries = [...]string{
	"daijirin",
	"daijisen",
	"kotowaza",
	"meikyou",
}

func JJLookupCommand(event *irc.Event) {
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
					result, err := findInDictionary(tokens[1], dict)

					// Unexpected errors
					if err != nil {
						fmt.Printf("ERROR: %s\n", err)

						failedLock.Lock()
						failed++
						failedLock.Unlock()

						wg.Done()
						return
					}

					// Not found
					if result == "" {
						failedLock.Lock()
						failed++
						failedLock.Unlock()

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

// findInDictionary searches for a given word inside a JSON dictionary file
// with structure `"word": "definition"`
func findInDictionary(word string, dictionary string) (string, error) {
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
