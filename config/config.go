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
	JwtKey             string
}

// NewConfig load env file and return a Config
func NewConfig(loadEnv bool) (Config, error) {
	var twitchClientId, twitchClientSecret, dbName, key string
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

	if key = os.Getenv("JWT_SIGNING_KEY"); key == "" {
		return Config{}, errors.New("JWT_SIGNING_KEY not set")
	}

	return Config{
		TwitchClientId:     twitchClientId,
		TwitchClientSecret: twitchClientSecret,
		DbName:             dbName,
		JwtKey:             key,
	}, nil
}
