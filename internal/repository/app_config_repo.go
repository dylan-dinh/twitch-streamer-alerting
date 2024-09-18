package repository

import (
	"errors"
	"fmt"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

type AppConfig interface {
	Insert(domain.AppConfig) error
	Get() (domain.AppConfig, error)
	Update(domain.AppConfig) error
}

type AppConfigRepo struct {
	Db     *gorm.DB
	logger *slog.Logger
}

func NewAppConfigRepo(db *gorm.DB) *AppConfigRepo {
	handler := slog.NewTextHandler(os.Stdout, nil)
	return &AppConfigRepo{Db: db, logger: slog.New(handler)}
}

func (a *AppConfigRepo) Insert(config domain.AppConfig) error {
	res := a.Db.Create(&config)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (a *AppConfigRepo) Get() (domain.AppConfig, error) {
	var config domain.AppConfig
	res := a.Db.First(&config)
	if res.Error != nil {
		return domain.AppConfig{}, res.Error
	}

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return domain.AppConfig{}, fmt.Errorf("record not found : %w", res.Error)
	}

	return config, nil
}

func (a *AppConfigRepo) Update(config domain.AppConfig) error {
	res := a.Db.Save(&config)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
