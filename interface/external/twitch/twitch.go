package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dylan-dinh/twitch-streamer-alerting/config"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/service"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const (
	twitchUrlAccessToken  = "https://id.twitch.tv/oauth2/token"
	grantType             = "client_credentials"
	baseDelay             = 1 * time.Second  // Initial delay between retries
	maxRetries            = 5                // Maximum number of retries
	maxDelay              = 64 * time.Second // Maximum backoff delay to avoid indefinite waiting
	refreshInterval       = 5 * time.Minute  // Interval to check if token needs refreshing
	JobRefreshAccessToken = "twitch job refresh access token"
)

type Twitch struct {
	logger             *slog.Logger
	http               *http.Client
	twitchClientId     string
	twitchClientSecret string
	BroadcasterService *service.BroadcasterService
	AppConfigService   *service.AppConfigService
}

func New(config config.Config,
	broadcasterService *service.BroadcasterService,
	appConfigService *service.AppConfigService) *Twitch {

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
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
		BroadcasterService: broadcasterService,
		AppConfigService:   appConfigService,
	}
}

func (t *Twitch) BackgroundUpdateBroadcasterInfo() {

}

func (t *Twitch) GetBroadcastersID() error {
	_, err := t.RefreshAccessToken()
	if err != nil {
		return err
	}

	token, err := t.AppConfigService.AppConfig.Get()
	if err != nil {
		return err
	}

	broadcasters, err := t.BroadcasterService.Broadcaster.GetBroadcastersWithoutUrl()
	if err != nil {
		return err
	}

	var userResponses = make([]UserResponse, len(broadcasters))
	for _, broadcaster := range broadcasters {
		uri := fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", broadcaster.Login)

		request, err := http.NewRequest(http.MethodGet, uri, nil)
		if err != nil {
			return err
		}
		request.Header.Add("Client-Id", t.twitchClientId)
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

		req, err := t.http.Do(request)
		if err != nil {
			return err
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			req.Body.Close()
			return err
		}

		var userResponse UserResponse
		if err := json.Unmarshal(body, &userResponse); err != nil {
			return err
		}

		if len(userResponse.Data) == 0 {
			return fmt.Errorf("user %s not found with url : %s", broadcaster.Login, uri)
		}

		userResponses = append(userResponses, userResponse)

		req.Body.Close()
	}
	return nil
}

func (t *Twitch) GetAccessToken() (ResponseAccessToken, error) {
	values := url.Values{}
	values.Set("client_id", t.twitchClientId)
	values.Set("client_secret", t.twitchClientSecret)
	values.Set("grantType", grantType)

	resp, err := t.http.PostForm(twitchUrlAccessToken, values)
	if err != nil {
		return ResponseAccessToken{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ResponseAccessToken{}, fmt.Errorf("error while getting access token, status code is : %d", resp.StatusCode)
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseAccessToken{}, err
	}

	var responseAccessToken ResponseAccessToken
	if err = json.Unmarshal(all, &responseAccessToken); err != nil {
		return ResponseAccessToken{}, err
	}

	return responseAccessToken, nil
}

// RefreshAccessToken check in db if access token has expired
// if expired it refreshes the access token
// if not expired we get it from db
func (t *Twitch) RefreshAccessToken() (bool, error) {
	toRefresh := false
	var token ResponseAccessToken
	conf, err := t.AppConfigService.AppConfig.Get()
	if err != nil {
		return false, err
	}

	compare := time.Now().Compare(conf.ExpiresAt)
	if compare == 0 || compare == +1 {
		token, err = t.GetAccessToken()
		if err != nil {
			t.logger.Error(err.Error())
			return false, err
		}

		err = t.AppConfigService.AppConfig.Update(domain.AppConfig{
			AccessToken: token.AccessToken,
			ExpiresAt:   time.Now().Add(time.Second * time.Duration(token.ExpiresIn)),
		})
		if err != nil {
			return false, err
		}
		toRefresh = true
		return toRefresh, nil
	}

	return toRefresh, nil
}

// retryExponentialBackoff attempts to retry a function with exponential backoff
func (t *Twitch) retryExponentialBackoff() bool {
	for attempt := 0; attempt < maxRetries; attempt++ {
		toRefresh, err := t.RefreshAccessToken()
		if toRefresh {
			// Calculate exponential backoff : bitwise left shift to double at every attempt * random
			// integer needed
			delayMillis := int64(baseDelay/time.Millisecond) * (1 << attempt)
			if delayMillis > int64(maxDelay/time.Millisecond) {
				delayMillis = int64(maxDelay / time.Millisecond) // Cap the delay to maxDelay
			}

			// Add a random factor (0.5 + rand.Float64()) to the delay
			waitTime := time.Duration(float64(delayMillis)*(0.5+rand.Float64())) * time.Millisecond
			t.logger.Info("Attempt %d failed: %v. Retrying in %v secs...\n", attempt+1, err, waitTime)
			time.Sleep(waitTime)
		} else {
			return true // Token refreshed successfully
		}
	}
	return false // Max retries reached without success
}

// BackgroundRefreshAccessToken runs the token refresh in the background and returns an error channel
func (t *Twitch) BackgroundRefreshAccessToken(ctx context.Context, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	go func() {
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				slog.Log(ctx, slog.LevelInfo, "stopping twitch job to refresh access token")
				return
			case <-ticker.C:
				if !t.retryExponentialBackoff() {
					errChan <- fmt.Errorf("critical: failed to refresh access token after max retries")
					return
				} else {
					t.logger.Info("Access token refreshed successfully")
				}
			}
		}
	}()
}
