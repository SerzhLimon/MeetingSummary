package service

import (
	"strconv"
	"strings"
)

func (w *Worker) SaveIncomingVoice(voiceBytes []byte, formatAudio string, chatID int64) (int, error) {
	return w.storage.SaveIncomingVoice(voiceBytes, formatAudio, chatID)
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

func (w *Worker) GigaChatReqBuilder(args []string) string {
	var builder strings.Builder
	for i := range args {
		builder.WriteString(args[i])
		builder.WriteString(" ")
	}

	return builder.String()
}

func (w *Worker) FindByKeyWords(chatID int64, keyWords []string) ([]int64, error) {

	patterns := make([]string, len(keyWords))
	for i := range keyWords {
		var builder strings.Builder
		builder.WriteString("%")
		builder.WriteString(keyWords[i])
		builder.WriteString("%")
		patterns[i] = builder.String()
	}
	return w.storage.FindByKeyWords(chatID, patterns)
}