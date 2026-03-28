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
        err := rows.Scan(&uploadData.VoiceData)
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

func (s *Storage) SetStatusFail(voiceID int) error {
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

func (s *Storage) SetStatusUpload(voiceID int, reqFileID string) error {
    _, err := s.db.Exec(querySetStatusUpload, voiceID, models.Upload, reqFileID)
    return err
}