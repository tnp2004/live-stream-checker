package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Twitch *Twitch
}

type Twitch struct {
	ClientID     string
	ClientSecret string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	return &Config{
		Twitch: &Twitch{
			ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		},
	}
}
