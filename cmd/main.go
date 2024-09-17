package main

import (
	"github.com/dylan-dinh/twitch-streamer-alerting/config"
	"github.com/dylan-dinh/twitch-streamer-alerting/interface/db"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
)

func main() {
	newConfig, err := config.NewConfig(true)
	if err != nil {
		panic(err)
	}

	//t := twitch.New(newConfig)
	//token, err := t.GetAccessToken()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(token)

	sqlite, err := db.NewSqlite(newConfig)
	if err != nil {
		panic(err)
	}

	email := "test@gmail.com"
	sqlite.Db.Create(&domain.User{Email: &email})
}
