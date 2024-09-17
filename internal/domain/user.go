package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email       *string
	FirstName   *string
	LastName    *string
	Broadcaster []*Broadcaster `gorm:"many2many:user_broadcasters;"`
}

type Broadcaster struct {
	gorm.Model
	BroadcasterId uint8
}
