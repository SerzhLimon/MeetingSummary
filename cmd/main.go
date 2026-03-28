package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Инициализация БД
	dbClient, err := db.InitPostgresClient(&cfg.Postgres)
	if err != nil {
		logrus.Fatalln(err)
	}
	
	err = migrations.Up(dbClient)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer func() {
		migrations.Down(dbClient)
		logrus.Info("Migrations down")
	}()
	storage := storage.New(dbClient)

	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	saluteWorker := service.InitWorker(cfg, storage)
	go saluteWorker.Run(ctx)

	bot := telebot.New(cfg, storage)
	bot.Route()
	
	go func() {
		logrus.Info("Starting bot...")
		bot.Start()
	}()
	
	// Ожидаем сигналы завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down gracefully...")
	
	cancel()
	
	shutdownTimeout := 2 * time.Second
	time.Sleep(shutdownTimeout)
	
	bot.Stop()
	
	logrus.Info("Shutdown completed")
}