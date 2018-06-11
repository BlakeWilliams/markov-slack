package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/nlopes/slack"
)

func main() {
	fmt.Println("Connecting to Slack...")
	bot_api_key, present := os.LookupEnv("BOT_API_KEY")
	if !present {
		fmt.Println("BOT_API_KEY environtment variable missing")
		os.Exit(1)
	}

	api := slack.New(bot_api_key)

	var bot_id string
	var bot_name string

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			bot_id = ev.Info.User.ID
			bot_name = ev.Info.User.Name
		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

			bot_id_regex := regexp.MustCompile(
				fmt.Sprintf("(?i)<@%s>|%s|business", bot_id, bot_name),
			)

			if bot_id_regex.MatchString(ev.Text) {
				fmt.Println("YOOOO!")
			}
		}
	}
}
