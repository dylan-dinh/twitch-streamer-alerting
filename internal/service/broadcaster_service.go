package service

import "github.com/dylan-dinh/twitch-streamer-alerting/internal/repository"

type BroadcasterService struct {
	Broadcaster repository.Broadcaster
}

func NewBroadcasterService(broadcaster repository.Broadcaster) *BroadcasterService {
	return &BroadcasterService{
		Broadcaster: broadcaster,
	}
}
