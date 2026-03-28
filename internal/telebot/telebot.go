package telebot

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/SerzhLimon/MeetingSummary/internal/config"
	s "github.com/SerzhLimon/MeetingSummary/internal/storage"
	"github.com/sirupsen/logrus"
	tg "gopkg.in/telebot.v3"
)

type TeleBot struct {
	core    *tg.Bot
	storage *s.Storage
}

func New(cfg *config.Config, storage *s.Storage) *TeleBot {
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

	return &TeleBot{
		core:    bot,
		storage: storage,
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
			logrus.Error(err)
			return c.Send(errSaveVoice)
		}
		voiceBytes, err := io.ReadAll(file)
		if err != nil {
			logrus.Error(err)
			return c.Send(errSaveVoice)
		}
		id, err := b.storage.SaveIncommingVoice(voiceBytes)
		if err != nil {
			logrus.Error(err)
			return c.Send(errSaveVoice)
		}

		return c.Send(fmt.Sprintf(successSaveVoice, id))
	})
}
