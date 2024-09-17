package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	TwitchClientId     string
	TwitchClientSecret string
	DbName             string
}

// NewConfig load env file and return a Config
func NewConfig(loadEnv bool) (Config, error) {
	var twitchClientId, twitchClientSecret, dbName string
	if loadEnv {
		err := godotenv.Load()
		if err != nil {
			return Config{}, err
		}
	}

	if twitchClientId = os.Getenv("TWITCH_CLIENT_ID"); twitchClientId == "" {
		return Config{}, errors.New("TWITCH_CLIENT_ID not set")
	}

	if twitchClientSecret = os.Getenv("TWITCH_CLIENT_SECRET"); twitchClientSecret == "" {
		return Config{}, errors.New("TWITCH_CLIENT_SECRET not set")
	}

	if dbName = os.Getenv("SQLITE_DB_NAME"); dbName == "" {
		return Config{}, errors.New("SQLITE_DB_NAME not set")
	}

	return Config{
		TwitchClientId:     twitchClientId,
		TwitchClientSecret: twitchClientSecret,
		DbName:             dbName,
	}, nil
}
