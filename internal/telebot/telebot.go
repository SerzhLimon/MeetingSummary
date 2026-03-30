package telebot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/SerzhLimon/MeetingSummary/internal/config"
	"github.com/SerzhLimon/MeetingSummary/internal/models"
	"github.com/SerzhLimon/MeetingSummary/internal/service"
	s "github.com/SerzhLimon/MeetingSummary/internal/storage"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	tg "gopkg.in/telebot.v3"
)

type TeleBot struct {
	core   *tg.Bot
	worker *service.Worker
	stopCh chan struct{}
}

func New(cfg *config.Config, storage *s.Storage) *TeleBot {
	settings := tg.Settings{
		Token:  cfg.Bot.Token,
		Poller: &tg.LongPoller{Timeout: 3 * time.Second},
	}

	bot, err := tg.NewBot(settings)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	w := service.InitWorker(cfg, storage)
	return &TeleBot{
		core:   bot,
		worker: w,
		stopCh: make(chan struct{}),
	}
}

func (b *TeleBot) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				logrus.Info("TeleBot: received shutdown signal, stopping...")
				return
			case m := <-b.worker.MessageChannel:
				b.core.Send(&telebot.Chat{ID: m.ChatID}, m.Message)
			}
		}
	}()

	go b.RunWorker(ctx)

	go func() {
		b.core.Start()
		close(b.stopCh)
	}()
}

func (b *TeleBot) Stop() {
	b.core.Stop()
	<-b.stopCh
}

func (b *TeleBot) Route() {
	b.core.Handle("/get", b.getHandler)
	b.core.Handle(tg.OnVoice, b.voiceHandler)
}

func (b *TeleBot) RunWorker(ctx context.Context) {
	go b.worker.Run(ctx)
}

func (b *TeleBot) voiceHandler(c tg.Context) error {
	msg := c.Message().Voice

	file, err := b.core.File(&tg.File{FileID: msg.FileID})
	if err != nil {
		logrus.Error(err)
		return c.Send(models.MsgErrSaveVoice)
	}
	voiceBytes, err := io.ReadAll(file)
	if err != nil {
		logrus.Error(err)
		return c.Send(models.MsgErrSaveVoice)
	}

	id, err := b.worker.SaveIncomingVoice(voiceBytes, c.Chat().ID)
	if err != nil {
		logrus.Error(err)
		return c.Send(models.MsgErrSaveVoice)
	}

	return c.Send(fmt.Sprintf(models.MsgSuccessSaveVoice, id))
}

func (b *TeleBot) getHandler(c tg.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send(models.GetErrEmptyID)
	}
	voiceID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return c.Send(models.GetErrInvalidID)
	}
	summary, err := b.worker.GetSummaryByID(voiceID, c.Chat().ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Send(models.Get404)
		}
		return c.Send(models.MsgInternalServerError)
	}
	return c.Send(summary)
}

// func (b *TeleBot) listHandler(c tg.Context) error {
// 	summary, err := b.worker.GetSummaryByID(voiceID, c.Chat().ID)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return c.Send(models.Get404)
// 		}
// 		return c.Send(models.MsgInternalServerError)
// 	}
// 	return c.Send(summary)
// }