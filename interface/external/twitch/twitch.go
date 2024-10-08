package twitch

import (
	"fmt"
	"github.com/dylan-dinh/twitch-streamer-alerting/config"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/service"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	twitchUrlAccessToken  = "https://id.twitch.tv/oauth2/token"
	grantType             = "client_credentials"
	baseDelay             = 1 * time.Second  // Initial delay between retries
	maxRetries            = 5                // Maximum number of retries
	maxDelay              = 64 * time.Second // Maximum backoff delay to avoid indefinite waiting
	refreshInterval       = 1 * time.Minute  // Interval to check if token needs refreshing
	JobRefreshAccessToken = "twitch job refresh access token"
)

type BroadcasterInfoResponse struct {
	Data []struct {
		ID          string `json:"id"`
		Login       string `json:"login"`
		DisplayName string `json:"display_name"`
		Type        string `json:"type"`
	} `json:"data"`
}

type ResponseAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (bir *BroadcasterInfoResponse) ToBroadcasterDomain() ([]domain.Broadcaster, error) {
	broadcasters := make([]domain.Broadcaster, len(bir.Data))

	for _, b := range bir.Data {
		atoi, err := strconv.ParseUint(b.ID, 10, 32)
		if err != nil {
			return nil, err
		}
		broadcasters = append(broadcasters, domain.Broadcaster{
			BroadcasterId: atoi,
			Login:         b.Login,
			DisplayName:   b.DisplayName,
			Type:          b.Type,
			Url:           fmt.Sprintf("https://twitch.tv/%s", b.Login),
		})
	}
	return broadcasters, nil
}

type Twitch struct {
	logger             *slog.Logger
	http               *http.Client
	twitchClientId     string
	twitchClientSecret string
	BroadcasterService *service.BroadcasterService
	AppConfigService   *service.AppConfigService
	IsTokenRefreshed   bool
}

func New(config config.Config,
	broadcasterService *service.BroadcasterService,
	appConfigService *service.AppConfigService) *Twitch {

	handler := slog.NewTextHandler(os.Stdout, nil)
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
