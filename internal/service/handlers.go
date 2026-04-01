package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SerzhLimon/MeetingSummary/internal/models"
)

func (w *Worker) uploadExecute(upload models.UploadData) (string, error) {
	w.getTokenSalute()

	reqBody := bytes.NewReader(upload.VoiceData)

	req, err := http.NewRequest("POST", "https://smartspeech.sber.ru/rest/v1/data:upload", reqBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+w.authSalute.AccessToken)
	req.Header.Set("Content-Type", upload.Format)

	resp, err := w.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("updload handler return status %d", resp.StatusCode)
	}

	response := models.UploadResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Result.RequestFileID, nil
}

func (w *Worker) recognizeExecute(recognize models.RecognizeData) (string, error) {
	w.getTokenSalute()

	var audioEncoding string
	switch recognize.Format {
	case models.FormatOgg:
		audioEncoding = "OPUS"
	case models.FormatMp3:
		audioEncoding = "MP3"
	}

	requestBody := models.RecognizeRequest{
		Options: models.RecognizeRequestOptions{
			Model:         "general",
			AudioEncoding: audioEncoding,
			SampleRate:    16000,
			ChannelsCount: 1,
		},
		RequestFileID: recognize.ReqFileID,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://smartspeech.sber.ru/rest/v1/speech:async_recognize", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+w.authSalute.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("recognize handler return status %d", resp.StatusCode)
	}

	response := models.RecognizeResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Result.ID, nil
}

func (w *Worker) checkStatusExecute(checkStatus models.CheckStatusData) (string, error) {
	w.getTokenSalute()

	url := fmt.Sprintf("https://smartspeech.sber.ru/rest/v1/task:get?id=%s", checkStatus.RecognizeID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+w.authSalute.AccessToken)

	resp, err := w.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("check status handler return status %d", resp.StatusCode)
	}

	response := models.ResponseCheckStatus{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	if response.Result.Status != "DONE" {
		return "", models.VoiceIsProseccing
	}
	return response.Result.ResponseFileID, nil
}

func (w *Worker) downloadTranscriptionExecute(download models.DownloadTranscriptionData) (string, error) {
	w.getTokenSalute()

	url := fmt.Sprintf("https://smartspeech.sber.ru/rest/v1/data:download?response_file_id=%s", download.RespFileID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+w.authSalute.AccessToken)

	resp, err := w.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download handler return status %d", resp.StatusCode)
	}

	response := []models.ResponseDownloadData{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
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

func (w *Worker) createSummaryExecute(summaryData models.CreateSummaryData) (string, error) {

	requestBody := models.GigaChatRequest{
		Model:             "GigaChat",
		Stream:            false,
		RepetitionPenalty: 1,
		Messages: []models.Message{
			{
				Role:    "user",
				Content: fmt.Sprintf("сделай краткую выжимку из текста: %s\n результат должен быть информативным и без лишней воды", summaryData.Text),
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://gigachat.devices.sberbank.ru/api/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+w.authGigaChat.AccessToken)

	resp, err := w.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("create summary handler return status %d", resp.StatusCode)
	}

	response := models.GigaChatResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	if len(response.Choices) < 1 {
		return "", fmt.Errorf("create summary handler return status empty response")
	}
	var summary string
	for i := range response.Choices {
		summary += response.Choices[i].Message.Content
	}

	return summary, nil
}
