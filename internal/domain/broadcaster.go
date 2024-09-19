package domain

import "gorm.io/gorm"

type Broadcaster struct {
	gorm.Model
	BroadcasterId uint8
	Login         string
	DisplayName   string
	Type          string
	Url           string
}
