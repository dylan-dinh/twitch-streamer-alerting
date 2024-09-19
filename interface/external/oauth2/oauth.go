package oauth2

import (
	"github.com/dylan-dinh/twitch-streamer-alerting/interface/external/twitch"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/service"
	"log/slog"
	"os"
)

type OauthClient interface {
	GetAccessToken() (twitch.ResponseAccessToken, error)
	RefreshAccessToken() (bool, error)
}

type OauthService struct {
	appConfig *service.AppConfigService
	logger    *slog.Logger
}

func NewOauthService(appConfig *service.AppConfigService) *OauthService {
	handler := slog.NewTextHandler(os.Stdout, nil)
	return &OauthService{
		logger:    slog.New(handler),
		appConfig: appConfig,
	}
}
