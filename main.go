package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tnp2004/live-stream-checker/checker"
)

type Config struct {
	twitch Twitch
}

type Twitch struct {
	ClientID     string
	ClientSecret string
}

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

func loadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		twitch: Twitch{
			ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		},
	}
}
