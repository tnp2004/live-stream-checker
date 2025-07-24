package checker

import (
	"log"
	"sync"

	"github.com/tnp2004/live-stream-checker/config"
	"github.com/tnp2004/live-stream-checker/models"
)

const (
	YOUTUBE = "youtube"
	TWITCH  = "twitch"
)

var (
	checkerInstance Checker
)

type Checker struct {
	youtube *youtube
	twitch  *twitch
}

type IChecker interface {
	IsLive(string) (bool, error)
}

func New(ch *models.Channel, cfg *config.Config) IChecker {
	var once sync.Once

	switch ch.Platform {
	case YOUTUBE:
		once.Do(func() {
			checkerInstance.youtube = NewYoutube()
			log.Println("Created youtube instance")
		})
		return checkerInstance.youtube
	case TWITCH:
		once.Do(func() {
			checkerInstance.twitch = NewTwitch(cfg.Twitch.ClientID, cfg.Twitch.ClientSecret)
			log.Println("Created twitch instance")
		})
		return checkerInstance.twitch
	}

	return nil
}
