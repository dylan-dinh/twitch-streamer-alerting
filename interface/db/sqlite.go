package db

import (
	"fmt"
	"github.com/dylan-dinh/twitch-streamer-alerting/config"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Sqlite struct {
	Db *gorm.DB
}

func NewSqlite(config config.Config) (*Sqlite, error) {
	db, err := gorm.Open(sqlite.Open(config.DbName), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(&domain.User{}, &domain.Broadcaster{}, &domain.AppConfig{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) GetDB() *gorm.DB {
	return s.Db
}
