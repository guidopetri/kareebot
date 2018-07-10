package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/thoj/go-ircevent"
)

var lock sync.Mutex
var doodPath = "./doods.txt"

func DoodCommand(event *irc.Event) {
	go func(event *irc.Event) {
		// Split message into tokens
		tokens := strings.Fields(event.Message())
		channel := event.Arguments[0]

		if channel != "#rtk" {
			return
		}

		if len(tokens) == 0 {
			irccon.Privmsg("#rtk", "credits go to karageko for reporting a bug that crashes kareebot in this scenario")
			return
		}

		if tokens[0] == "!dood" {
			irccon.Privmsg("#rtk", fmt.Sprintf("dood counter: %d", getDoodCounter()))
			return
		}
		
		if tokens[0] == "!baka" {
			irccon.Privmsg("#rtk", fmt.Sprintf("karageko b-baka!"))
			return
		}

		raw := strings.ToLower(event.Message())
		count := strings.Count(raw, "dood")

		if count > 0 {
			fmt.Printf("dooded: %d", count)
			incrementDoodCounter(count)
		}

	}(event)
}

func incrementDoodCounter(addedDoods int) bool {
	count := getDoodCounter()

	lock.Lock()
	defer lock.Unlock()

	f, err := os.Create(doodPath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer f.Close()

	io.WriteString(f, fmt.Sprintf("%d", count+addedDoods))

	return true
}

func getDoodCounter() int {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.Open(doodPath)
	if err != nil {
		fmt.Println(err)
		os.Create(doodPath)
		return 0
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	c, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return c
}
