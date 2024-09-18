package repository

import (
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

type Broadcaster interface {
	Insert(broadcaster domain.Broadcaster) error
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
