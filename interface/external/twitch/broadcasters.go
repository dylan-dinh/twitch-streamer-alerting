package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func (t *Twitch) BackgroundUpdateBroadcasterInfo(ctx context.Context, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	go func() {
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				slog.Log(ctx, slog.LevelInfo, "stopping update of twitch broadcaster info")
				return
			case <-ticker.C:
				//TODO: add the func to get and update broadcaster info
			}
		}
	}()
}

func (t *Twitch) GetBroadcastersID() ([]BroadcasterInfoResponse, error) {
	_, err := t.RefreshAccessToken()
	if err != nil {
		return nil, err
	}

	token, err := t.AppConfigService.AppConfig.Get()
	if err != nil {
		return nil, err
	}

	broadcasters, err := t.BroadcasterService.Broadcaster.GetBroadcastersWithoutUrl()
	if err != nil {
		return nil, err
	}
	values := url.Values{}

	var userResponses = make([]BroadcasterInfoResponse, 0)
	for _, broadcaster := range broadcasters {
		values.Add("login", broadcaster.Login)
	}

	uri := fmt.Sprintf("https://api.twitch.tv/helix/users?%s", values.Encode())

	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Client-Id", t.twitchClientId)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	req, err := t.http.Do(request)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var userResponse BroadcasterInfoResponse
	if err = json.Unmarshal(body, &userResponse); err != nil {
		return nil, err
	}

	if len(userResponse.Data) == 0 {
		t.logger.Info("no users updated from twitch api")
	}

	userResponses = append(userResponses, userResponse)

	return userResponses, nil
}
