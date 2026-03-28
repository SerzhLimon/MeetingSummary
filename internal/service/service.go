package service

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/SerzhLimon/MeetingSummary/internal/config"
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
		}
	}
}

func (w *SaluteWorker) upload() {
	uploadData, err := w.storage.GetVoiceForUpload()
	if err != nil {
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
			err = w.storage.SetStatusFail(uploadData[i].VoiceID)
			//
		}
	}
}
