package main

import (
	"fmt"
	"github.com/dylan-dinh/twitch-streamer-alerting/api"
	"github.com/dylan-dinh/twitch-streamer-alerting/config"
	"github.com/dylan-dinh/twitch-streamer-alerting/interface/db"
	"github.com/dylan-dinh/twitch-streamer-alerting/interface/external/oauth2"
	"github.com/dylan-dinh/twitch-streamer-alerting/interface/external/twitch"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/factory"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/jwt"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/repository"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/service"
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
	DB := sqlite.GetDB()

	// internal repositories
	appConfigRepo := repository.NewAppConfigRepo(DB)
	broadcasterRepo := repository.NewBroadcasterRepo(DB)
	userRepo := repository.NewUserRepo(DB)

	// internal services
	appConfigService := service.NewAppconfigService(appConfigRepo)
	broadcastService := service.NewBroadcasterService(broadcasterRepo)

	// oauth services
	_ = oauth2.NewOauthService(appConfigService)

	// external services
	twitchService := twitch.New(newConfig, broadcastService, appConfigService)

	_, err = twitchService.GetBroadcastersID()
	if err != nil {
		panic(err)
	}

	// Jwt service
	j := jwt.NewJwt(newConfig)

	// handler initialization
	user := api.NewUserHandler(service.NewUserService(userRepo, DB), j)

	// router set up
	router := api.SetUpRouter(user)

	// routine setup
	routines := []factory.Routines{
		{
			Name:    twitch.JobRefreshAccessToken,
			Routine: twitchService.BackgroundRefreshAccessToken,
		},
	}

	routinesFactory := factory.NewRoutinesFactory(routines)
	routinesFactory.StartRoutinesFactory()

	// Run the HTTP server
	go func() {
		if err := router.Run(":8080"); err != nil {
			routinesFactory.ErrChan <- fmt.Errorf("gin server error: %v", err)
		}
	}()

	routinesFactory.StopRoutinesFactory()
}
