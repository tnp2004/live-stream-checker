package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

const FIND_CHANNEL_ID_REGEX = `https:\/\/www\.youtube\.com\/channel\/([a-zA-Z0-9_-]+)`
const CHECK_LIVE_REGEX = `<link\s+rel=["']canonical["']\s+href=["']([^"']+)["']`

type youtubeResponse struct {
	Items []struct {
		Snippet struct {
			LiveBroadcastContent string `json:"liveBroadcastContent"`
		} `json:"snippet"`
	} `json:"items"`
}

type youtube struct {
	apiKey string
}

func NewYoutube(googleApiKey string) *youtube {
	return &youtube{googleApiKey}
}

func (yt youtube) IsLive(url string) (bool, error) {
	pageSource, err := yt.getPageSource(url)
	if err != nil {
		return false, err
	}
	channelID, err := yt.findChannelID(pageSource)
	if err != nil {
		return false, err
	}

	return yt.checkLive(channelID)
}

func (yt youtube) getPageSource(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error: ", err.Error())
		return "", fmt.Errorf("get page source from %s failed", url)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error: ", err.Error())
		return "", fmt.Errorf("read response body failed")
	}
	pageSource := string(body)

	return pageSource, nil
}

func (yt youtube) findChannelID(pageSource string) (string, error) {
	regex := regexp.MustCompile(FIND_CHANNEL_ID_REGEX)
	match := regex.FindStringSubmatch(pageSource)
	if match == nil {
		log.Println("Error: channel id not found")
		return "", fmt.Errorf("channel id not found")
	}
	channelID := match[1]

	return channelID, nil
}

func (yt youtube) checkLive(channelID string) (bool, error) {
	url := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=%s&type=video&eventType=live&key=%s",
		channelID, yt.apiKey,
	)
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error: ", err.Error())
		return false, fmt.Errorf("get data failed")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error: ", err.Error())
		return false, fmt.Errorf("read response body failed")
	}
	respBody := new(youtubeResponse)
	log.Println(respBody)
	if err := json.Unmarshal(body, &respBody); err != nil {
		log.Println("Error: ", err.Error())
		return false, err
	}
	if len(respBody.Items) == 0 || respBody.Items[0].Snippet.LiveBroadcastContent != "live" {
		return false, nil
	}

	return true, nil
}
