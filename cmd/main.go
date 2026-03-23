package main

import (
	"github.com/SerzhLimon/MeetingSummary/internal/config"
	"github.com/SerzhLimon/MeetingSummary/internal/telebot"
)

func main() {
	cfg := config.LoadConfig()
	bot := telebot.New(cfg)
	bot.Start()
	defer bot.Stop()
}
