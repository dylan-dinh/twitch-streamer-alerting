package repository

import (
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

type Broadcaster interface {
	Insert(domain.Broadcaster) error
	GetBroadcastersWithoutUrl() ([]domain.Broadcaster, error)
}

type BroadcasterRepo struct {
	Db     *gorm.DB
	logger *slog.Logger
}

func NewBroadcasterRepo(db *gorm.DB) *BroadcasterRepo {
	handler := slog.NewTextHandler(os.Stdout, nil)
	return &BroadcasterRepo{Db: db, logger: slog.New(handler)}
}

func (b *BroadcasterRepo) Insert(broadcaster domain.Broadcaster) error {
	res := b.Db.Create(&broadcaster)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (b *BroadcasterRepo) GetBroadcastersWithoutUrl() ([]domain.Broadcaster, error) {
	var broadcasters []domain.Broadcaster
	res := b.Db.Where("url IS NULL").Find(&broadcasters)
	if res.Error != nil {
		return []domain.Broadcaster{}, res.Error
	}
	return broadcasters, nil
}

func (b *BroadcasterRepo) Update(broadcaster domain.Broadcaster) error {
	res := b.Db.Save(&broadcaster)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
