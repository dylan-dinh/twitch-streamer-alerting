package domain

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrEmailExists = errors.New("email already exists")
)

// BadRequestError custom context error type
type BadRequestError struct {
	Err error
}

func (e *BadRequestError) Error() string {
	return e.Err.Error()
}

type User struct {
	gorm.Model
	ID          string  `gorm:"type:text;primaryKey"`
	Email       *string `gorm:"unique;not null"`
	Username    *string `gorm:"not null"`
	FirstName   *string
	LastName    *string
	Password    *string        `gorm:"not null"`
	Broadcaster []*Broadcaster `gorm:"many2many:user_broadcasters;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String() // Generate a new UUID v4 as a string
	return
}
