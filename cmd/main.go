package main

import (
	"github.com/dylan-dinh/twitch-streamer-alerting/config"
	"github.com/dylan-dinh/twitch-streamer-alerting/interface/db"
	"github.com/dylan-dinh/twitch-streamer-alerting/interface/external/oauth2"
	"github.com/dylan-dinh/twitch-streamer-alerting/interface/external/twitch"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/repository"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/service"
	"log"
)

func main() {
	// config to load env var
	newConfig, err := config.NewConfig(true)
	if err != nil {
		panic(err)
	}

	// sqlite connection
	sqlite, err := db.NewSqlite(newConfig)
	if err != nil {
		panic(err)
	}

	DB := sqlite.Db

	// external services
	twitchService := twitch.New(newConfig)

	// internal repositories
	appConfigRepo := repository.NewAppConfigRepo(DB)
	broadcasterRepo := repository.NewBroadcasterRepo(DB)

	// internal services
	appConfigService := service.New(appConfigRepo)
	_ = service.NewBroadcasterService(broadcasterRepo)

	// oauth services
	oauthService := oauth2.NewOauthService(twitchService, appConfigService)

	errChan := oauthService.BackgroundRefreshAccessToken()
	for {
		select {
		case err = <-errChan:
			log.Fatalf("error from background process: %v", err)
		}
	}
}
