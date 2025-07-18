package checker

import (
	"strings"
)

type twitch struct {
	url         string
	accessToken string
	channelName string
}

func Twitch(url string) twitch {
	twitch := twitch{url: url}
	channelName := twitch.getTwitchChannelName()
	twitch.channelName = channelName
	return twitch
}

func (tw twitch) IsLive() (bool, error) {
	return false, nil
}

func (tw twitch) requestTwitchAccessToken() string {
	return ""
}

func (tw twitch) getTwitchChannelName() string {
	twitchUrlPrefix := "https://www.twitch.tv/"
	channelName := strings.TrimPrefix(tw.url, twitchUrlPrefix)
	slashIndex := strings.Index(channelName, "/")
	if slashIndex != -1 {
		channelName = channelName[:slashIndex]
	}
	return channelName
}
