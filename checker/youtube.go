package checker

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const FIND_CHANNEL_ID_REGEX = `https:\/\/www\.youtube\.com\/channel\/([a-zA-Z0-9_-]+)`
const CHECK_LIVE_REGEX = `<link\s+rel=["']canonical["']\s+href=["']([^"']+)["']`

type youtube struct {
	url string
}

func Youtube(url string) youtube {
	return youtube{url}
}

func (yt youtube) IsLive() (bool, error) {
	pageSource, err := yt.getPageSource(yt.url)
	if err != nil {
		return false, err
	}
	channelID, err := yt.findChannelID(pageSource)
	if err != nil {
		return false, err
	}
	liveUrl := fmt.Sprintf("https://www.youtube.com/channel/%s/live", channelID)
	pageSource, err = yt.getPageSource(liveUrl)
	if err != nil {
		return false, err
	}

	return yt.checkLive(pageSource)
}

func (yt youtube) getPageSource(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return "", fmt.Errorf("get page source from %s failed", url)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return "", fmt.Errorf("read response body failed")
	}
	pageSource := string(body)

	return pageSource, nil
}

func (yt youtube) findChannelID(pageSource string) (string, error) {
	regex := regexp.MustCompile(FIND_CHANNEL_ID_REGEX)
	match := regex.FindStringSubmatch(pageSource)
	if match == nil {
		fmt.Println("Error: channel id not found")
		return "", fmt.Errorf("channel id not found")
	}
	channelID := match[1]

	return channelID, nil
}

func (yt youtube) checkLive(pageSource string) (bool, error) {
	regex := regexp.MustCompile(CHECK_LIVE_REGEX)
	match := regex.FindStringSubmatch(pageSource)
	if match == nil {
		fmt.Println("Error: canonical not found")
		return false, fmt.Errorf("check Live failed")
	}
	liveUrl := match[1]
	if !strings.Contains(liveUrl, "watch") {
		return false, nil
	}

	return true, nil
}
