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

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (w *Worker) getTokenSalute() {
	if w.authSalute.ExpiresAt.After(time.Now()) {
		return
	}
	rquid := uuid.New().String()
	body := fmt.Sprintf("scope=%s", "SALUTE_SPEECH_PERS")
	reqBody := strings.NewReader(body)

	req, err := http.NewRequest("POST", "https://ngw.devices.sberbank.ru:9443/api/v2/oauth", reqBody)
	if err != nil {
		logrus.Error("Worker.getTokenSalute(): ", err)
		return
	}

	authKey := base64.StdEncoding.EncodeToString([]byte(w.cfg.Salute.ClientID + ":" + w.cfg.Salute.ClientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("RqUID", rquid)
	req.Header.Set("Authorization", "Basic "+authKey)

	resp, err := w.client.Do(req)
	if err != nil {
		logrus.Error("Worker.getTokenSalute(): ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error("Worker.getTokenSalute(): ", resp.StatusCode)
		return
	}

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		logrus.Error("Worker.getTokenSalute(): ", err)
	}
	token := result["access_token"]
	expiresAt := result["expires_at"]

	w.authSalute.AccessToken = token.(string)
	w.authSalute.ExpiresAt = time.Unix(int64(expiresAt.(float64)), 0)
}

func (w *Worker) getTokenGigaChat() {
	if w.authGigaChat.ExpiresAt.After(time.Now()) {
		return
	}

	rquid := uuid.New().String()
	authKey := base64.StdEncoding.EncodeToString([]byte(w.cfg.GigaChat.ClientID + ":" + w.cfg.GigaChat.ClientSecret))
	body := []byte("scope=GIGACHAT_API_PERS")

	req, err := http.NewRequest("POST", "https://ngw.devices.sberbank.ru:9443/api/v2/oauth", bytes.NewBuffer(body))
	if err != nil {
		logrus.Error("Worker.getTokenGigaChat(): ", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", rquid)
	req.Header.Set("Authorization", "Basic "+authKey)

	resp, err := w.client.Do(req)
	if err != nil {
		logrus.Error("Worker.getTokenGigaChat(): ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error("Worker.getTokenGigaChat(): ", resp.StatusCode)
		return
	}

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		logrus.Error("Worker.getTokenGigaChat(): ", err)
	}
	token := result["access_token"]
	expiresAt := result["expires_at"]

	w.authGigaChat.AccessToken = token.(string)
	w.authGigaChat.ExpiresAt = time.Unix(int64(expiresAt.(float64)), 0)
}
