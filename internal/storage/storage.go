package storage

import "database/sql"

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db:db}
}

func (s *Storage) SaveIncommingVoice(voiceBytes []byte) (int, error) {
	var voiceID int
	err := s.db.QueryRow(string(voiceBytes)).Scan(&voiceID)
	return voiceID, err
}