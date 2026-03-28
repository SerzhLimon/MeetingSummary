package storage

import (
	"database/sql"
	"fmt"

	"github.com/SerzhLimon/MeetingSummary/internal/models"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db:db}
}

func (s *Storage) SaveIncomingVoice(voiceBytes []byte) (int, error) {
	var voiceID int
	err := s.db.QueryRow(querySaveIncomingVoice, voiceBytes, models.Begin).Scan(&voiceID)
	return voiceID, err
}

func (s *Storage) GetVoiceForUpload() ([]models.UploadData, error) {
    rows, err := s.db.Query(queryGetVoicesForUpload, models.Begin)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var uploadDataList []models.UploadData
    
    for rows.Next() {
        var uploadData models.UploadData
        err := rows.Scan(&uploadData.VoiceID, &uploadData.VoiceData)
        if err != nil {
            return nil, err
        }
        uploadDataList = append(uploadDataList, uploadData)
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }
    
    return uploadDataList, nil
}

func (s *Storage) SetStatusFail(voiceID int64) error {
    result, err := s.db.Exec(querySetStatusFail, voiceID, models.Fail)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("storage.SetStatusFail(): rows affected = 0 for voice: %d", voiceID)
	}
    return nil
}

func (s *Storage) SetStatusUpload(voiceID int64, reqFileID string) error {
    _, err := s.db.Exec(querySetStatusUpload, voiceID, models.Upload, reqFileID)
    return err
}

func (s *Storage) GetVoiceForRecognize() ([]models.RecognizeData, error) {
    rows, err := s.db.Query(queryGetVoicesForRecognize, models.Upload)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var recognizeDataList []models.RecognizeData
    
    for rows.Next() {
        var recognizeData models.RecognizeData
        err := rows.Scan(&recognizeData.VoiceID, &recognizeData.ReqFileID)
        if err != nil {
            return nil, err
        }
        recognizeDataList = append(recognizeDataList, recognizeData)
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }
    
    return recognizeDataList, nil
}

func (s *Storage) SetStatusRecognition(voiceID int64, recognizeID string) error {
    _, err := s.db.Exec(querySetStatusRecognition, voiceID, models.Recognition, recognizeID)
    return err
}