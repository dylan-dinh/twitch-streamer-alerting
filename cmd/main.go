package main

import (
	"fmt"
	"github.com/dylan-dinh/twitch-streamer-alerting/config"
	"github.com/dylan-dinh/twitch-streamer-alerting/interface/db"
	"github.com/dylan-dinh/twitch-streamer-alerting/interface/external/oauth2"
	"github.com/dylan-dinh/twitch-streamer-alerting/interface/external/twitch"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/factory"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/repository"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Set up routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

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

	// internal repositories
	appConfigRepo := repository.NewAppConfigRepo(DB)
	broadcasterRepo := repository.NewBroadcasterRepo(DB)

	// internal services
	appConfigService := service.NewAppconfigService(appConfigRepo)
	broadcastService := service.NewBroadcasterService(broadcasterRepo)

	// oauth services
	_ = oauth2.NewOauthService(appConfigService)

	// external services
	twitchService := twitch.New(newConfig, broadcastService, appConfigService)

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
		if err := r.Run(":8080"); err != nil {
			routinesFactory.ErrChan <- fmt.Errorf("gin server error: %v", err)
		}
	}()

	routinesFactory.StopRoutinesFactory()

}
