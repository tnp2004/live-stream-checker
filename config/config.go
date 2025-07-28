package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Google *Google
	Twitch *Twitch
}

type Twitch struct {
	ClientID     string
	ClientSecret string
}

type Google struct {
	ApiKey string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	return &Config{
		Google: &Google{ApiKey: os.Getenv("GOOGLE_API_KEY")},
		Twitch: &Twitch{
			ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		},
	}
}
