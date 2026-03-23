package telebot

import (
	"log"
	"time"

	"github.com/SerzhLimon/MeetingSummary/internal/config"
	tg "gopkg.in/telebot.v3"
)

type TeleBot struct {
	core *tg.Bot
}

func New(cfg *config.Config) *TeleBot {
	settings := tg.Settings{
		Token:  cfg.Bot.Token,
		Poller: &tg.LongPoller{Timeout: 3 * time.Second},
	}

	// Создаём бота
	bot, err := tg.NewBot(settings)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	bot.Handle("/start", func(c tg.Context) error {
		return c.Send("Привет, бродяга")
	})
	
	return &TeleBot{core: bot}
}
func (b *TeleBot) Start() {
	log.Println("success start bot")
	b.core.Start()
}

func (b *TeleBot) Stop() {
	b.core.Stop()
}

// func (b *TeleBot) Route() {
// 	b.core.Handle("/start", func(c tg.Context) error {
// 		return c.Send("Привет")
// 	})
	
// }
