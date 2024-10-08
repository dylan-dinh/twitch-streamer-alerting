package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email       *string `gorm:"unique;not null"`
	Username    *string `gorm:"not null"`
	FirstName   *string
	LastName    *string
	Password    *string        `gorm:"not null"`
	Broadcaster []*Broadcaster `gorm:"many2many:user_broadcasters;"`
}
