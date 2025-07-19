package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tnp2004/live-stream-checker/checker"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	checker.Twitch("https://www.twitch.tv/charmer_cham/clip/GiantIntelligentKiwiDatSheffy-kHLWQbe2iYBLsug3").IsLive()
}
