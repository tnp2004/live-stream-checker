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
	once            sync.Once
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
	once.Do(func() {
		checkerInstance.youtube = NewYoutube(cfg.Google.ApiKey)
		checkerInstance.twitch = NewTwitch(cfg.Twitch.ClientID, cfg.Twitch.ClientSecret)
		log.Println("Created instance")
	})

	switch ch.Platform {
	case YOUTUBE:
		return checkerInstance.youtube
	case TWITCH:
		return checkerInstance.twitch
	}

	return nil
}
