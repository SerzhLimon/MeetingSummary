package main

import (
	"github.com/SerzhLimon/MeetingSummary/internal/config"
	"github.com/SerzhLimon/MeetingSummary/internal/config/db"
	"github.com/SerzhLimon/MeetingSummary/internal/service"
	"github.com/SerzhLimon/MeetingSummary/internal/storage"
	"github.com/SerzhLimon/MeetingSummary/internal/telebot"
	"github.com/SerzhLimon/MeetingSummary/migrations"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()

	db, err := db.InitPostgresClient(&cfg.Postgres)
	if err != nil {
		logrus.Fatalln(err)
	}
	err = migrations.Up(db)
	if err != nil {
		logrus.Fatalln(err)
	}
	storage := storage.New(db)

	saluteWorker := service.InitWorker(cfg, storage)

	bot := telebot.New(cfg, storage)
	bot.Route()
	bot.Start()
	defer bot.Stop()
}
