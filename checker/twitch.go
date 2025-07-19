package checker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const TWITCH_REQUEST_TOKEN_URL = "https://id.twitch.tv/oauth2/token"
const TWITCH_GRANT_TYPE = "client_credentials"
const LIVE_CHECK_URL = "https://api.twitch.tv/helix/streams?user_login="

type twitch struct {
	accessToken  string
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

type liveStatus struct {
	Data []struct {
		Username string `json:"user_name"`
		Title    string `json:"title"`
		GameName string `json:"game_name"`
		Type     string `json:"type"`
	} `json:"data"`
}

func NewTwitch(clientID, clientSecret string) *twitch {
	twitch := &twitch{clientID: clientID, clientSecret: clientSecret}
	twitch.getAccessToken()
	return twitch
}

func (tw twitch) IsLive(url string) (bool, error) {
	channelName := getChannelName(url)
	return tw.checkLive(channelName)
}

func (tw twitch) checkLive(channelName string) (bool, error) {
	url := fmt.Sprintf("%s%s", LIVE_CHECK_URL, channelName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		return false, err
	}
	bearerToken := fmt.Sprintf("Bearer %s", tw.accessToken)
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("Client-ID", tw.clientID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
		return false, err
	}
	if resp.StatusCode != 200 {
		fmt.Printf("Error: http status %s", resp.Status)
		return false, fmt.Errorf("checking live error")
	}
	defer resp.Body.Close()

	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return false, err
	}
	respBody := new(liveStatus)
	if err := json.Unmarshal(byteBody, &respBody); err != nil {
		fmt.Println("Error: ", err.Error())
		return false, err
	}

	if len(respBody.Data) == 0 {
		return false, nil
	}

	return true, nil
}

func (tw *twitch) getAccessToken() error {
	if len(tw.accessToken) != 0 {
		return nil
	}
	accessToken, err := requestAccessToken(tw.clientID, tw.clientSecret, TWITCH_GRANT_TYPE)
	if err != nil {
		return err
	}
	tw.accessToken = accessToken

	return nil
}

func requestAccessToken(clientID, clientSecret, grantType string) (string, error) {
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

func getChannelName(url string) string {
	twitchUrlPrefix := "https://www.twitch.tv/"
	channelName := strings.TrimPrefix(url, twitchUrlPrefix)
	slashIndex := strings.Index(channelName, "/")
	if slashIndex != -1 {
		channelName = channelName[:slashIndex]
	}
	return channelName
}
