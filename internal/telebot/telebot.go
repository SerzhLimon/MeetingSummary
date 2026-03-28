package telebot

import (
	"fmt"
	"log"
	"time"

	"github.com/SerzhLimon/MeetingSummary/internal/config"
	tg "gopkg.in/telebot.v3"
)

type TeleBot struct {
	core *tg.Bot
	// uc   *s.Service

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
	// uc := s.New()
	return &TeleBot{
		core: bot,
		// uc:   uc,
	}
}

func (b *TeleBot) Start() {
	log.Println("success start bot")
	b.core.Start()
}

func (b *TeleBot) Stop() {
	b.core.Stop()
}

func (b *TeleBot) Route() {
	b.core.Handle("/start", func(c tg.Context) error {
		return c.Send("Привет, бродяга")
	})
	b.core.Handle(tg.OnVoice, func(c tg.Context) error {
		msg := c.Message().Voice

		file, err := b.core.File(&tg.File{FileID: msg.FileID})
		if err != nil {
			return c.Send("Ошибка загрузки файла")
		}
		fmt.Print(file)
		// err = b.uc.SendToRemoteServer(file)
		if err != nil {
			return c.Send("Ошибка отправки на сервер")
		}

		return nil
	})
}
