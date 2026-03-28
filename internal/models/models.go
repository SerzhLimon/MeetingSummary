package models

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	VoiceIsProseccing Err = "встреча еще в обработке"


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

type CheckStatusData struct {
	VoiceID     int64
	RecognizeID string
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

type resultCheckStatus struct {
	ID             string `json:"id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	Status         string `json:"status"`
	ResponseFileID string `json:"response_file_id"`
}
type ResponseCheckStatus struct {
	Status int               `json:"status"`
	Result resultCheckStatus `json:"result"`
}
