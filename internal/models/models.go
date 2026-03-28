package models

const (
	Fail        = "FAIL"
	Begin       = "BEGIN"
	Upload      = "UPLOAD"
	Recognition = "RECOGNITION"
	Wait        = "WAIT"
	Download    = "DOWNLOAD"
)

type UploadData struct {
	VoiceID   int64
	VoiceData []byte
}

type RecognizeData struct {
	VoiceID   int64
	ReqFileID string
}

type resultUpload struct {
	RequestFileID string `json:"request_file_id"`
}

type UploadResponse struct {
	Result resultUpload `json:"result"`
}

type resultRecognize struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Status    string `json:"status"`
}
type RecognizeResponse struct {
	Status int             `json:"status"`
	Result resultRecognize `json:"result"`
}
