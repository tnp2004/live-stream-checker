package main

import (
	"fmt"

	"github.com/tnp2004/live-stream-checker/checker"
)

func main() {
	config := loadConfig()
	twitch := checker.NewTwitch(config.twitch.ClientID, config.twitch.ClientSecret)
	url := "https://www.twitch.tv/takluz"
	isLive, _ := twitch.IsLive(url)
	if !isLive {
		fmt.Println("Not live")
	} else {
		fmt.Println("Live")
	}
}
