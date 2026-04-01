package models

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	VoiceIsProseccing  Err = "meeting processed yet"
	NoDataForProcessed Err = "no voices for processed"

	Fail        = "FAIL"
	Begin       = "BEGIN"
	Upload      = "UPLOAD"
	Recognition = "RECOGNITION"
	Wait        = "WAIT"
	Download    = "DOWNLOAD"
	Success     = "SUCCESS"

	FormatOgg = "audio/ogg"
	FormatMp3 = "audio/mpeg"
)

type UserMessage struct {
	ChatID  int64
	Message string
}

type UploadData struct {
	VoiceID   int64
	ChatID    int64
	Format    string
	VoiceData []byte
}

type RecognizeData struct {
	VoiceID   int64
	ChatID    int64
	Format    string
	ReqFileID string
}

type CheckStatusData struct {
	VoiceID     int64
	ChatID      int64
	RecognizeID string
}

type DownloadTranscriptionData struct {
	VoiceID    int64
	ChatID     int64
	RespFileID string
}

type CreateSummaryData struct {
	VoiceID int64
	ChatID  int64
	Text    string
}

type resultUpload struct {
	RequestFileID string `json:"request_file_id"`
}

type UploadResponse struct {
	Result resultUpload `json:"result"`
}

type RecognizeRequestOptions struct {
	Model         string `json:"model"`
	AudioEncoding string `json:"audio_encoding"`
	SampleRate    int    `json:"sample_rate"`
	ChannelsCount int    `json:"channels_count"`
}

type RecognizeRequest struct {
	Options       RecognizeRequestOptions `json:"options"`
	RequestFileID string                  `json:"request_file_id"`
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

type wordAlignment struct {
	Word  string `json:"word"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type result struct {
	Text           string          `json:"text"`
	NormalizedText string          `json:"normalized_text"`
	Start          string          `json:"start"`
	End            string          `json:"end"`
	WordAlignments []wordAlignment `json:"word_alignments"`
}

type emotionsResult struct {
	Positive float64 `json:"positive"`
	Neutral  float64 `json:"neutral"`
	Negative float64 `json:"negative"`
}

type backendInfo struct {
	ModelName     string `json:"model_name"`
	ModelVersion  string `json:"model_version"`
	ServerVersion string `json:"server_version"`
}

type speakerInfo struct {
	SpeakerID             int     `json:"speaker_id"`
	MainSpeakerConfidence float64 `json:"main_speaker_confidence"`
}

type personIdentity struct {
	Age         string `json:"age"`
	Gender      string `json:"gender"`
	AgeScore    int    `json:"age_score"`
	GenderScore int    `json:"gender_score"`
}

type ResponseDownloadData struct {
	Results             []result       `json:"results"`
	Eou                 bool           `json:"eou"`
	EmotionsResult      emotionsResult `json:"emotions_result"`
	ProcessedAudioStart string         `json:"processed_audio_start"`
	ProcessedAudioEnd   string         `json:"processed_audio_end"`
	BackendInfo         backendInfo    `json:"backend_info"`
	Channel             int            `json:"channel"`
	SpeakerInfo         speakerInfo    `json:"speaker_info"`
	EouReason           string         `json:"eou_reason"`
	Insight             string         `json:"insight"`
	PersonIdentity      personIdentity `json:"person_identity"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GigaChatRequest struct {
	Model             string    `json:"model"`
	Messages          []Message `json:"messages"`
	Stream            bool      `json:"stream"`
	RepetitionPenalty int       `json:"repetition_penalty"`
}

type choicesResponse struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
}

type usageResponse struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	SystemTokens     int `json:"system_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type GigaChatResponse struct {
	Choices []choicesResponse `json:"choices"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Object  string            `json:"object"`
	Usage   usageResponse     `json:"usage"`
}
