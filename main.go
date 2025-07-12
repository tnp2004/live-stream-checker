package main

import (
	"log"

	"github.com/tnp2004/live-stream-checker/checker"
)

func main() {
	youtube := checker.NewYoutube("https://www.youtube.com/@TEDx")
	isLive, err := youtube.IsLive()
	if err != nil {
		return
	}
	if !isLive {
		log.Println("Not Live . . .")
		return
	}
	log.Println("Live . . .")
}
