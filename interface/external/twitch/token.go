package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

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
	if compare == 0 || compare == 1 {
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
		t.IsTokenRefreshed = toRefresh
		return toRefresh, nil
	}

	t.IsTokenRefreshed = toRefresh
	return toRefresh, nil
}

// retryExponentialBackoff attempts to retry a function with exponential backoff
func (t *Twitch) retryExponentialBackoff() bool {
	for attempt := 0; attempt < maxRetries; attempt++ {
		toRefresh, err := t.RefreshAccessToken()
		if toRefresh {
			// Calculate exponential backoff : bitwise left shift to double at every attempt * random
			// int64 needed for comparison and operation
			delayMillis := int64(baseDelay/time.Millisecond) * (1 << attempt)
			if delayMillis > int64(maxDelay/time.Millisecond) {
				delayMillis = int64(maxDelay / time.Millisecond) // Cap the delay to maxDelay
			}

			// Add a random factor (0.5 + rand.Float64()) to the delay
			waitTime := time.Duration(float64(delayMillis)*(0.5+rand.Float64())) * time.Millisecond
			t.logger.Info(fmt.Sprintf("Attempt %d failed: %v. Retrying in %v secs...\n", attempt+1, err, waitTime))
			time.Sleep(waitTime)
		} else {
			return true // Token refreshed successfully
		}
	}
	return false // Max retries reached without success
}

// BackgroundRefreshAccessToken runs the token refresh in the background
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
					if t.IsTokenRefreshed {
						t.logger.Info("Access token refreshed successfully")
					}
					t.logger.Info("Access token was not expired yet")

				}
			}
		}
	}()
}
