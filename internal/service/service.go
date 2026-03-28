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
	}

	w.getTokenSalute()
	w.getTokenGigaChat()
	return w
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(w.cfg.IntervalTicker) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Info("Worker: received shutdown signal, stopping...")
			return
		case <-ticker.C:
			// logic
			w.upload()
			w.recognize()
			w.checkStatus()
			w.downloadTranscription()
			w.createSummary()
		}
	}
}

func (w *Worker) upload() {
	uploadData, err := w.storage.GetVoiceForUpload()
	if err != nil {
		if errors.Is(err, models.NoDataForProcessed) {
			logrus.Warn(err)
			return
		}
		logrus.Error(err)
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
		} else {
			logrus.Infof("success upload: %d", uploadData[i].VoiceID)
		}
	}
}

func (w *Worker) recognize() {
	recognizeData, err := w.storage.GetVoiceForRecognize()
	if err != nil {
		//
		logrus.Error(fmt.Errorf("Worker.recognize(): %w", err))
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
		} else {
			logrus.Infof("success recognize: %d", recognizeData[i].VoiceID)
		}
	}
}

func (w *Worker) checkStatus() {
	checkStatusData, err := w.storage.GetVoiceForCheckStatus()
	if err != nil {
		//
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
			err = w.storage.SetStatusFail(checkStatusData[i].VoiceID)
			//
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
		//
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
			err = w.storage.SetStatusFail(downloadTranscriptionData[i].VoiceID)
			//
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
		//
		logrus.Error(fmt.Errorf("Worker.downloadTranscription(): %w", err))
		return
	}
	for i := range createSumData {
		summary, err := w.createSummaryExecute(createSumData[i])
		if err != nil {
			logrus.Error(err)
			err = w.storage.SetStatusFail(createSumData[i].VoiceID)
			//
			continue
		}
		// err = w.storage.SetStatusDownload(downloadTranscriptionData[i].VoiceID, text)
		// if err != nil {
		// 	logrus.Error(err)
		// 	continue
		// }
		logrus.Infof("successfull create summary: %d", createSumData[i].VoiceID)
		logrus.Info("DONE--------\n",summary)
	}
}
