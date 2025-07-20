package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	twitch *Twitch
}

type Twitch struct {
	ClientID     string
	ClientSecret string
}

func loadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		twitch: &Twitch{
			ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		},
	}
}
