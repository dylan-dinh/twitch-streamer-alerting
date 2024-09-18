package oauth2

import (
	"fmt"
	"github.com/dylan-dinh/twitch-streamer-alerting/interface/external/twitch"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/service"
	"log/slog"
	"math/rand"
	"os"
	"time"
)

const (
	baseDelay       = 1 * time.Second  // Initial delay between retries
	maxRetries      = 5                // Maximum number of retries
	maxDelay        = 64 * time.Second // Maximum backoff delay to avoid indefinite waiting
	refreshInterval = 5 * time.Minute  // Interval to check if token needs refreshing
)

type OauthClient interface {
	GetAccessToken() (twitch.ResponseAccessToken, error)
}

type OauthService struct {
	client    OauthClient
	appConfig service.AppConfigService
	logger    *slog.Logger
}

func NewOauthService(client OauthClient, appConfig service.AppConfigService) *OauthService {
	handler := slog.NewTextHandler(os.Stdout, nil)
	return &OauthService{
		logger:    slog.New(handler),
		client:    client,
		appConfig: appConfig,
	}
}

// RefreshAccessToken check in db if access token has expired
// if expired it refreshes the access token
// if not expired we get it from db
func (o *OauthService) RefreshAccessToken() (bool, error) {
	toRefresh := false
	var token twitch.ResponseAccessToken
	conf, err := o.appConfig.AppConfig.Get()
	if err != nil {
		return false, err
	}

	compare := time.Now().Compare(conf.ExpiresAt)
	if compare == 0 || compare == +1 {
		token, err = o.client.GetAccessToken()
		if err != nil {
			o.logger.Error(err.Error())
			return false, err
		}

		err = o.appConfig.AppConfig.Update(domain.AppConfig{
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
func (o *OauthService) retryExponentialBackoff() bool {
	for attempt := 0; attempt < maxRetries; attempt++ {
		toRefresh, err := o.RefreshAccessToken()
		if toRefresh {
			// Calculate exponential backoff : bitwise left shift to double at every attempt * random
			waitTime := time.Duration(float64(baseDelay) * (1 << attempt) * (0.5 + rand.Float64()))
			if waitTime > maxDelay {
				waitTime = maxDelay // Cap the delay to maxDelay
			}
			o.logger.Info("Attempt %d failed: %v. Retrying in %v secs...\n", attempt+1, err, waitTime)
			time.Sleep(waitTime)
		} else {
			return true // Token refreshed successfully
		}
	}
	return false // Max retries reached without success
}

// BackgroundRefreshAccessToken runs the token refresh in the background and returns an error channel
func (o *OauthService) BackgroundRefreshAccessToken() chan error {
	errChan := make(chan error)

	go func() {
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				success := o.retryExponentialBackoff()
				if !success {
					errChan <- fmt.Errorf("critical: failed to refresh access token after max retries")
					return
				} else {
					o.logger.Info("Access token refreshed successfully")
				}
			}
		}
	}()
	return errChan
}
