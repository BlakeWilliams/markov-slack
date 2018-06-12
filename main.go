package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

func main() {
	markov := NewMarkov()

	go fetchHistory(markov)
	runRtm(markov)
}

func fetchHistory(markov *MarkovChain) {
	apiKey, present := os.LookupEnv("USER_API_KEY")
	if !present {
		fmt.Println("USER_API_KEY environtment variable missing")
		os.Exit(1)
	}

	userName, present := os.LookupEnv("USER_NAME")
	if !present {
		fmt.Println("USER_NAME environtment variable missing")
		os.Exit(1)
	}

	api := slack.New(apiKey)

	params := slack.NewSearchParameters()
	params.Count = 100

	searchQuery := fmt.Sprintf("from:@%s", userName)

	for {
		messages, _, error := api.Search(searchQuery, params)
		if error != nil {
			fmt.Println("error searching", error)
			thirtySeconds := time.Duration(30) * time.Second
			time.Sleep(thirtySeconds)

			continue
		}

		for _, message := range messages.Matches {
			// ensures message isn't private channel or private chat
			if strings.HasPrefix(message.Channel.ID, "C") {
				markov.Parse(message.Text)
			}
		}

		if messages.Paging.Page == messages.Paging.Pages {
			break
		} else {
			params.Page = messages.Paging.Page + 1
		}
	}
}

func runRtm(markov *MarkovChain) {
	apiKey, present := os.LookupEnv("BOT_API_KEY")
	if !present {
		fmt.Println("BOT_API_KEY environtment variable missing")
		os.Exit(1)
	}

	api := slack.New(apiKey)

	var botId string
	var botName string

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			botId = ev.Info.User.ID
			botName = ev.Info.User.Name
		case *slack.MessageEvent:
			botIdRegex := regexp.MustCompile(
				fmt.Sprintf("(?i)<@%s>|%s|business", botId, botName),
			)

			if botIdRegex.MatchString(ev.Text) {
				message := markov.GenerateSentence()

				rtm.SendMessage(
					rtm.NewOutgoingMessage(message, ev.Channel),
				)
			}
		}
	}
}
