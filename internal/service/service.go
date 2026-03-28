package service

import (
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/SerzhLimon/MeetingSummary/internal/config"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Auth struct {
	AccessToken string
	ExpiresAt   time.Time
}

type SaluteWorker struct {
	cfg     *config.Config
	storage *sql.DB
	client  *http.Client
	auth    *Auth
}

func InitWorker(cfg *config.Config, db *sql.DB) *SaluteWorker {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return &SaluteWorker{
		cfg:     cfg,
		storage: db,
		client:  client,
	}
}

func (w *SaluteWorker) Auth() {
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
