package domain

import "gorm.io/gorm"

type Broadcaster struct {
	gorm.Model
	BroadcasterId uint8
	Name          string
	Url           string
}
