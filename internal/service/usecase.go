package service

func (w *Worker) SaveIncomingVoice(voiceBytes []byte, chatID int64) (int, error) {
	return w.storage.SaveIncomingVoice(voiceBytes, chatID)
}

func (w *Worker) GetSummaryByID(voiceID, chatID int64) (string, error) {
	return w.storage.GetSummaryByID(voiceID, chatID)
}