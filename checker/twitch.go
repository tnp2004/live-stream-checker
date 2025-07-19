package checker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const TWITCH_REQUEST_TOKEN_URL = "https://id.twitch.tv/oauth2/token"
const TWITCH_GRANT_TYPE = "client_credentials"

type twitch struct {
	url          string
	accessToken  string
	channelName  string
	clientID     string
	clientSecret string
}

type accessTokenRequestBody struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type accessTokenResponseBody struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"bearer"`
}

func Twitch(url string) *twitch {
	channelName := getTwitchChannelName(url)
	return &twitch{
		url:          url,
		channelName:  channelName,
		clientID:     os.Getenv("TWITCH_CLIENT_ID"),
		clientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
	}
}

func (tw twitch) IsLive() (bool, error) {
	accessToken := tw.getTwitchAccessToken()
	return false, nil
}

func (tw twitch) getTwitchAccessToken() string {
	if len(tw.accessToken) != 0 {
		return tw.accessToken
	}
	accessToken, err := requestTwitchAccessToken(tw.clientID, tw.clientSecret, TWITCH_GRANT_TYPE)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return ""
	}

	return accessToken
}

func requestTwitchAccessToken(clientID, clientSecret, grantType string) (string, error) {
	body, err := json.Marshal(accessTokenRequestBody{clientID, clientSecret, grantType})
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return "", err
	}
	resp, err := http.Post(TWITCH_REQUEST_TOKEN_URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return "", err
	}

	respBody := new(accessTokenResponseBody)
	if err := json.Unmarshal(byteBody, &respBody); err != nil {
		fmt.Println("Error: ", err.Error())
		return "", err
	}

	return respBody.AccessToken, nil
}

func getTwitchChannelName(url string) string {
	twitchUrlPrefix := "https://www.twitch.tv/"
	channelName := strings.TrimPrefix(url, twitchUrlPrefix)
	slashIndex := strings.Index(channelName, "/")
	if slashIndex != -1 {
		channelName = channelName[:slashIndex]
	}
	return channelName
}
