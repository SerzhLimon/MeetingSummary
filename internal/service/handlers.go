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

	response := models.UploadResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.Error("SaluteWorker.uploadExecute(): ", err)
		return "", err
	}
	
	return response.Result.RequestFileID, nil
}

func (w *SaluteWorker) recognizeExecute(recognize models.RecognizeData) (string, error) {
// must return recognize_id
	w.getToken()
	requestBody := map[string]interface{}{
        "options":         map[string]interface{}{},
        "model":           "general",
        "audio_encoding":  "OPUS",
        "sample_rate":     16000,
        "channels_count":  1,
        "request_file_id": recognize.ReqFileID,
    }
    
    jsonBody, err := json.Marshal(requestBody)
    if err != nil {
        logrus.Error("SaluteWorker.recognizeExecute(): ", err)
        return "", err
    }

	req, err := http.NewRequest("POST", "https://smartspeech.sber.ru/rest/v1/speech:async_recognize", bytes.NewBuffer(jsonBody))
	if err != nil {
		logrus.Error("SaluteWorker.recognizeExecute(): ", err)
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+w.auth.AccessToken)
    req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		logrus.Error("SaluteWorker.recognizeExecute(): ", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error("SaluteWorker.recognizeExecute(): ", err)
		return "", err
	}

	response := models.RecognizeResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        logrus.Error("SaluteWorker.recognizeExecute(): ", err)
        return "", err
    }

	return response.Result.ID, nil
}

func (w *SaluteWorker) checkStatusExecute(checkStatus models.CheckStatusData) (string, error) {
	w.getToken()

	url := fmt.Sprintf("https://smartspeech.sber.ru/rest/v1/task:get?id=%s", checkStatus.RecognizeID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error("SaluteWorker.checkStatusExecute(): ", err)
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+w.auth.AccessToken)

	resp, err := w.client.Do(req)
	if err != nil {
		logrus.Error("SaluteWorker.checkStatusExecute(): ", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error("SaluteWorker.checkStatusExecute(): ", err)
		return "", err
	}

	response := models.ResponseCheckStatus{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        logrus.Error("SaluteWorker.checkStatusExecute(): ", err)
        return "", err
    }
	if response.Result.Status != "DONE" {
		return "", models.VoiceIsProseccing
	}
	return response.Result.ID, nil
}

func (w *SaluteWorker) downloadTranscriptionExecute(download models.DownloadTranscriptionData) (string, error) {
    w.getToken()

    url := fmt.Sprintf("https://smartspeech.sber.ru/rest/v1/data:download?response_file_id=%s", download.RespFileID)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        logrus.Error("SaluteWorker.downloadTranscriptionExecute(): ", err)
        return "", err
    }
    req.Header.Set("Authorization", "Bearer "+w.auth.AccessToken)

    resp, err := w.client.Do(req)
    if err != nil {
        logrus.Error("SaluteWorker.downloadTranscriptionExecute(): ", err)
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        logrus.Errorf("SaluteWorker.downloadTranscriptionExecute(): status %d, body: %s", resp.StatusCode, string(body))
        return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
    }

    response := []models.ResponseDownloadData{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        logrus.Error("SaluteWorker.downloadTranscriptionExecute(): ", err)
        return "", err
    }

    if len(response) == 0 {
        return "", fmt.Errorf("empty response from download endpoint")
    }

    var text string
    for i := range response {
        for j := range response[i].Results {
            if text != "" {
                text += " "
            }
            text += response[i].Results[j].Text
        }
    }

    return text, nil
}