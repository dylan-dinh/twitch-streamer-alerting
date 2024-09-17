package twitch

import (
	"encoding/json"
	"fmt"
	"github.com/dylan-dinh/twitch-streamer-alerting/config"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	twitchUrlAccessToken = "https://id.twitch.tv/oauth2/token"
	grant_type           = "client_credentials"
)

type Twitch struct {
	logger             *slog.Logger
	http               *http.Client
	twitchClientId     string
	twitchClientSecret string
}

func New(config config.Config) *Twitch {
	level := &slog.LevelVar{}
	level.Set(0)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       level,
		ReplaceAttr: nil,
	})

	return &Twitch{
		logger: slog.New(handler),
		http: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   5 * time.Second,
		},
		twitchClientId:     config.TwitchClientId,
		twitchClientSecret: config.TwitchClientSecret,
	}
}

func (t *Twitch) GetAccessToken() (string, error) {
	values := url.Values{}
	values.Set("client_id", t.twitchClientId)
	values.Set("client_secret", t.twitchClientSecret)
	values.Set("grant_type", grant_type)

	resp, err := t.http.PostForm(twitchUrlAccessToken, values)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error while getting access token, status code is : %d", resp.StatusCode)
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseAccessToken ResponseAccessToken
	if err = json.Unmarshal(all, &responseAccessToken); err != nil {
		return "", err
	}

	fmt.Println(responseAccessToken)

	return responseAccessToken.AccessToken, nil
}
