package service

import (
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/repository"
	"gorm.io/gorm"
)

type UserService struct {
	Repo repository.User
	Db   *gorm.DB
}

func NewUserService(user repository.User, db *gorm.DB) UserService {
	return UserService{
		Repo: user,
		Db:   db,
	}
}
