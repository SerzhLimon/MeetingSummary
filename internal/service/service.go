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

type SaluteWorker struct {
	cfg     *config.Config
	storage *s.Storage
	client  *http.Client
	auth    *Auth
}

func InitWorker(cfg *config.Config, storage *s.Storage) *SaluteWorker {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	sw := &SaluteWorker{
		cfg:     cfg,
		storage: storage,
		client:  client,
		auth:    &Auth{},
	}

	sw.getToken()
	return sw
}

func (w *SaluteWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(w.cfg.Salute.IntervalTicker) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// logic
			w.upload()
			w.recognize()
			w.checkStatus()
			w.downloadTranscription()
		}
	}
}

func (w *SaluteWorker) upload() {
	uploadData, err := w.storage.GetVoiceForUpload()
	if err != nil {
		//
		logrus.Error(fmt.Errorf("SaluteWorker.upload(): %w", err))
		return
	}

	for i := range uploadData {
		reqFileID, err := w.uploadExecute(uploadData[i])
		if err != nil {
			logrus.Error(err)
			err = w.storage.SetStatusFail(uploadData[i].VoiceID)
			//
			continue
		}
		err = w.storage.SetStatusUpload(uploadData[i].VoiceID, reqFileID)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (w *SaluteWorker) recognize() {
	recognizeData, err := w.storage.GetVoiceForRecognize()
	if err != nil {
		//
		logrus.Error(fmt.Errorf("SaluteWorker.recognize(): %w", err))
		return
	}
	for i := range recognizeData {
		recognizeID, err := w.recognizeExecute(recognizeData[i])
		if err != nil {
			logrus.Error(err)
			err = w.storage.SetStatusFail(recognizeData[i].VoiceID)
			//
			continue
		}
		err = w.storage.SetStatusRecognition(recognizeData[i].VoiceID, recognizeID)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (w *SaluteWorker) checkStatus() {
	checkStatusData, err := w.storage.GetVoiceForCheckStatus()
	if err != nil {
		//
		logrus.Error(fmt.Errorf("SaluteWorker.checkStatus(): %w", err))
		return
	}

	for i := range checkStatusData {
		respFileID, err := w.checkStatusExecute(checkStatusData[i])
		if err != nil {
			if errors.Is(err, models.VoiceIsProseccing) {
				continue
			}
			logrus.Error(err)
			err = w.storage.SetStatusFail(checkStatusData[i].VoiceID)
			//
			continue
		}
		err = w.storage.SetStatusWait(checkStatusData[i].VoiceID, respFileID)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (w *SaluteWorker) downloadTranscription() {
	downloadTranscriptionData, err := w.storage.GetVoiceForDownloadTranscription()
	if err != nil {
		//
		logrus.Error(fmt.Errorf("SaluteWorker.checkStatus(): %w", err))
		return
	}
	for i := range downloadTranscriptionData {
		text, err := w.downloadTranscriptionExecute(downloadTranscriptionData[i])
		if err != nil {
			if errors.Is(err, models.VoiceIsProseccing) {
				continue
			}
			logrus.Error(err)
			err = w.storage.SetStatusFail(downloadTranscriptionData[i].VoiceID)
			//
			continue
		}
		err = w.storage.SetStatusDownload(downloadTranscriptionData[i].VoiceID, text)
		if err != nil {
			logrus.Error(err)
		}
	}
}