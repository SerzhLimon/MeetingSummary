package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/SerzhLimon/MeetingSummary/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (w *SaluteWorker) getToken() {
	if w.auth.ExpiresAt.Before(time.Now()) {
		return
	}
	rquid := uuid.New().String()
	body := fmt.Sprintf("scope=%s", "SALUTE_SPEECH_PERS")
	reqBody := strings.NewReader(body)

	req, err := http.NewRequest("POST", "https://ngw.devices.sberbank.ru:9443/api/v2/oauth", reqBody)
	if err != nil {
		logrus.Error("SaluteWorker.Auth(): ", err)
		return
	}

	authKey := base64.StdEncoding.EncodeToString([]byte(w.cfg.Salute.ClientID + ":" + w.cfg.Salute.ClientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("RqUID", rquid)
	req.Header.Set("Authorization", "Basic "+authKey)

	resp, err := w.client.Do(req)
	if err != nil {
		logrus.Error("SaluteWorker.Auth(): ", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		logrus.Error("SaluteWorker.Auth(): ", err)
	}
	token := result["access_token"]
	expiresAt := result["expires_at"]

	w.auth.AccessToken = token.(string)
	w.auth.ExpiresAt = time.UnixMilli(expiresAt.(int64))
}

func (w *SaluteWorker) uploadExecute(upload models.UploadData) (string, error) {
	w.getToken()

	reqBody := bytes.NewReader(upload.VoiceData)

	req, err := http.NewRequest("POST", "https://smartspeech.sber.ru/rest/v1/data/upload", reqBody)
	if err != nil {
		logrus.Error("SaluteWorker.uploadExecute(): ", err)
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+w.auth.AccessToken)
	req.Header.Set("Content-Type", "audio/ogg")
	resp, err := w.client.Do(req)
	if err != nil {
		logrus.Error("SaluteWorker.uploadExecute(): ", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error("SaluteWorker.uploadExecute(): ", err)
		return "", err
	}

	var response struct {
		Result struct {
			RequestFileID string `json:"request_file_id"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.Error("SaluteWorker.uploadExecute(): ", err)
		return "", err
	}
	
	return response.Result.RequestFileID, nil
}
