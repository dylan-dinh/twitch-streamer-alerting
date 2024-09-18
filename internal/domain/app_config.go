package domain

import (
	"gorm.io/gorm"
	"time"
)

type AppConfig struct {
	gorm.Model
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
