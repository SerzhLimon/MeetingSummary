package models

const (
	Fail = "FAIL"
	Begin = "BEGIN"
	Upload = "UPLOAD"
	Recognition = "RECOGNITION"
	Wait = "WAIT"
	Download = "DOWNLOAD"
)


type UploadData struct {
	VoiceID int
	VoiceData []byte
}