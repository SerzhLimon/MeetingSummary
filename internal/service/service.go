// тут лежит воркер и его основной флоу
package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SerzhLimon/MeetingSummary/internal/config"
	"github.com/SerzhLimon/MeetingSummary/internal/models"
	s "github.com/SerzhLimon/MeetingSummary/internal/storage"
	"github.com/sirupsen/logrus"
)

type Auth struct {
	AccessToken string
	ExpiresAt   time.Time
}

type Worker struct {
	cfg          *config.Config
	storage      *s.Storage
	client       *http.Client
	authSalute   *Auth
	authGigaChat *Auth
	MessageChannel chan models.UserMessage
}

func InitWorker(cfg *config.Config, storage *s.Storage) *Worker {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	w := &Worker{
		cfg:          cfg,
		storage:      storage,
		client:       client,
		authSalute:   &Auth{},
		authGigaChat: &Auth{},
		MessageChannel: make(chan models.UserMessage),
	}

	w.getTokenSalute()
	w.getTokenGigaChat()
	return w
}

// краткая суть работы ядра воркера
// воркер запускает тикер (регулируется конфигом) который раз в N секунд запускает основной флоу
// извлекаются все загруженные аудюхи, по очереди прокидываются в салют/гигачат для преобразования в саммрари

// в случае факапа на каком то из этапов, аудюхе присваивается статус фейл и юзеру приходит сооющение об ошибке
// в случае успешного прохождения каждого из этапов у аудио меняется статус и загружаются доп данные в БД
// при получении саммари, аудио получается статус САКСЕСС и юзер получает уведомление об успехе
func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(w.cfg.IntervalTicker) * time.Second)
	defer ticker.Stop()
	defer close(w.MessageChannel)
	for {
		select {
		case <-ctx.Done():
			logrus.Info("Worker: received shutdown signal, stopping...")
			return
		case <-ticker.C:
			w.upload()
			w.recognize()
			w.checkStatus()
			w.downloadTranscription()
			w.createSummary()
		}
	}
}

// извлекает все аудюхи со статусом БЕГИН и по очереди отправляет в салют на обработку
// в случае факапа присвоит статус фейл аудюхе
// в случае успеха - загрузит нужные данные в БД и обновит статус
// (все нижеописанные методы раюотают аналогично)
func (w *Worker) upload() {
	uploadData, err := w.storage.GetVoiceForUpload()
	if err != nil {
		if errors.Is(err, models.NoDataForProcessed) {
			logrus.Warn(err)
			return
		}
		logrus.Error(fmt.Errorf("Worker.upload(): %w", err))
		return
	}

	for i := range uploadData {
		reqFileID, err := w.uploadExecute(uploadData[i])
		if err != nil {
			logrus.Error(err)

			go w.sendMsgFail(uploadData[i].ChatID, uploadData[i].VoiceID)

			if err = w.storage.SetStatusFail(uploadData[i].VoiceID); err != nil {
				logrus.Error(fmt.Errorf("Worker.upload(): %w", err))
			}
			continue
		}
		err = w.storage.SetStatusUpload(uploadData[i].VoiceID, reqFileID)
		if err != nil {
			logrus.Error(err)
		} else {
			logrus.Infof("success upload: %d, %d", uploadData[i].VoiceID, uploadData[i].ChatID)
		}
	}
}

func (w *Worker) recognize() {
	recognizeData, err := w.storage.GetVoiceForRecognize()
	if err != nil {
		logrus.Error(fmt.Errorf("Worker.recognize(): %w", err))
		return
	}

	for i := range recognizeData {
		recognizeID, err := w.recognizeExecute(recognizeData[i])
		if err != nil {
			logrus.Error(err)

			go w.sendMsgFail(recognizeData[i].ChatID, recognizeData[i].VoiceID)

			if err = w.storage.SetStatusFail(recognizeData[i].VoiceID); err != nil {
				logrus.Error(fmt.Errorf("Worker.recognize(): %w", err))
			}
			continue
		}
		err = w.storage.SetStatusRecognition(recognizeData[i].VoiceID, recognizeID)
		if err != nil {
			logrus.Error(err)
		} else {
			logrus.Infof("success recognize: %d", recognizeData[i].VoiceID)
		}
	}
}

func (w *Worker) checkStatus() {
	checkStatusData, err := w.storage.GetVoiceForCheckStatus()
	if err != nil {
		logrus.Error(fmt.Errorf("Worker.checkStatus(): %w", err))
		return
	}

	for i := range checkStatusData {
		respFileID, err := w.checkStatusExecute(checkStatusData[i])
		if err != nil {
			if errors.Is(err, models.VoiceIsProseccing) {
				logrus.Infof("status not DONE yet: %d", checkStatusData[i].VoiceID)
				continue
			}

			logrus.Error(err)

			go w.sendMsgFail(checkStatusData[i].ChatID, checkStatusData[i].VoiceID)

			if err = w.storage.SetStatusFail(checkStatusData[i].VoiceID); err != nil {
				logrus.Error(fmt.Errorf("Worker.checkStatus(): %w", err))
			}
			continue
		}
		err = w.storage.SetStatusWait(checkStatusData[i].VoiceID, respFileID)
		if err != nil {
			logrus.Error(err)
		} else {
			logrus.Infof("success check status: %d", checkStatusData[i].VoiceID)
		}
	}
}

func (w *Worker) downloadTranscription() {
	downloadTranscriptionData, err := w.storage.GetVoiceForDownloadTranscription()
	if err != nil {
		logrus.Error(fmt.Errorf("Worker.downloadTranscription(): %w", err))
		return
	}

	for i := range downloadTranscriptionData {
		text, err := w.downloadTranscriptionExecute(downloadTranscriptionData[i])
		if err != nil {
			if errors.Is(err, models.VoiceIsProseccing) {
				continue
			}

			logrus.Error(err)

			go w.sendMsgFail(downloadTranscriptionData[i].ChatID, downloadTranscriptionData[i].VoiceID)

			if err = w.storage.SetStatusFail(downloadTranscriptionData[i].VoiceID); err != nil {
				logrus.Error(fmt.Errorf("Worker.downloadTranscription(): %w", err))
			}
			continue
		}
		err = w.storage.SetStatusDownload(downloadTranscriptionData[i].VoiceID, text)
		if err != nil {
			logrus.Error(err)
			continue
		}
		logrus.Infof("successfull download: %d", downloadTranscriptionData[i].VoiceID)
	}
}

func (w *Worker) createSummary() {
	createSumData, err := w.storage.GetVoiceForCreateSummary()
	if err != nil {
		logrus.Error(fmt.Errorf("Worker.createSummary(): %w", err))
		return
	}

	for i := range createSumData {
		summary, err := w.createSummaryExecute(createSumData[i])
		if err != nil {
			logrus.Error(err)

			go w.sendMsgFail(createSumData[i].ChatID, createSumData[i].VoiceID)

			if err = w.storage.SetStatusFail(createSumData[i].VoiceID); err != nil {
				logrus.Error(fmt.Errorf("Worker.createSummary(): %w", err))
			}
			continue
		}
		err = w.storage.SetStatusSuccess(createSumData[i].VoiceID, summary)
		if err != nil {
			logrus.Error(err)
			continue
		}

		go w.sendMsgSuccess(createSumData[i].ChatID,createSumData[i].VoiceID)

		logrus.Infof("successfull create summary: %d", createSumData[i].VoiceID)
	}
}

// кидают сообщения в канал для уведомления юзера о результате работы
func (w *Worker) sendMsgSuccess(chatID, voiceID int64) {
	w.MessageChannel <- models.UserMessage{
		ChatID: chatID,
		Message: fmt.Sprintf(models.MsgSuccesSummaryProcess, voiceID),
	}
}

func (w *Worker) sendMsgFail(chatID, voiceID int64) {
	w.MessageChannel <- models.UserMessage{
		ChatID: chatID,
		Message: fmt.Sprintf(string(models.MsgFailSummaryProcess), voiceID),
	}
}