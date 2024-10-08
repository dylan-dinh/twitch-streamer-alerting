package repository

import (
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

type User interface {
	Insert(domain.User) error
}

type UserRepo struct {
	Db     *gorm.DB
	logger *slog.Logger
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return UserRepo{
		Db:     db,
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (r UserRepo) Insert(user domain.User) error {
	res := r.Db.Create(&user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r UserRepo) FindByEmailAndPassword(email, password string) (domain.User, error) {
	var user domain.User
	res := r.Db.Where("email = ? AND password = ?", email, password).First(&user)
	if res.Error != nil {
		return domain.User{}, res.Error
	}
	return user, nil
}
