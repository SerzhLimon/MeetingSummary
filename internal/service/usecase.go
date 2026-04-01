package service

import (
	"strconv"
	"strings"
)

func (w *Worker) SaveIncomingVoice(voiceBytes []byte, chatID int64) (int, error) {
	return w.storage.SaveIncomingVoice(voiceBytes, chatID)
}

func (w *Worker) GetSummaryByID(voiceID, chatID int64) (string, error) {
	return w.storage.GetSummaryByID(voiceID, chatID)
}

func (w *Worker) GetListSummaryID(chatID int64) ([]int64, error) {
	return w.storage.GetListSummaryID(chatID)
}

func (w *Worker) ListResponseBuilder(IDs []int64) string {
	
	var builder strings.Builder
	builder.WriteString("Вот список ваших встреч:\n")
	for i := range IDs {
		builder.WriteString("Встреча ")
		builder.WriteString(strconv.FormatInt(IDs[i], 10))
		builder.WriteString("\n")
	}

	return builder.String()
}
