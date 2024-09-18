package service

import (
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/repository"
)

type AppConfigService struct {
	AppConfig repository.AppConfig
}

func New(AppConfig repository.AppConfig) AppConfigService {
	return AppConfigService{
		AppConfig: AppConfig,
	}
}
