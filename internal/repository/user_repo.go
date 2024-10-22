package repository

import (
	"errors"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

type User interface {
	Insert(*gorm.DB, domain.User) (domain.User, error)
	FindByEmailAndPassword(string, string) (domain.User, error)
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

func (r UserRepo) Insert(tx *gorm.DB, user domain.User) (domain.User, error) {
	res := tx.Create(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return domain.User{}, &domain.BadRequestError{Err: domain.ErrEmailExists}
		}
		return domain.User{}, res.Error
	}

	return user, nil
}

func (r UserRepo) FindByEmailAndPassword(email, password string) (domain.User, error) {
	var user domain.User
	res := r.Db.Where("email = ? AND password = ?", email, password).First(&user)
	if res.Error != nil {
		return domain.User{}, res.Error
	}
	return user, nil
}
