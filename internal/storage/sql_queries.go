package storage

const (
	querySaveIncomingVoice = `
		INSERT INTO voice_recognize_req (voice_data, process_step) 
        VALUES ($1, $2)
		RETURNING id
	`

	queryGetVoicesForUpload = `
		SELECT id, voice_data
		FROM voice_recognize_req
		WHERE process_step = $1
	`

	querySetStatusFail = `
		UPDATE voice_recognize_req 
		SET process_step = $2
		WHERE id = $1
	`

	querySetStatusUpload = `
		UPDATE voice_recognize_req 
		SET process_step = $2, request_file_id = $3
		WHERE id = $1
	`

	queryGetVoicesForRecognize = `
		SELECT id, request_file_id
		FROM voice_recognize_req
		WHERE process_step = $1
	`

	querySetStatusRecognition = `
		UPDATE voice_recognize_req 
		SET process_step = $2, recognize_id = $3
		WHERE id = $1
	`

	queryGetVoiceForCheckStatus = `
		SELECT id, recognize_id
		FROM voice_recognize_req
		WHERE process_step = $1
	`

	querySetStatusWait = `
		UPDATE voice_recognize_req 
		SET process_step = $2, response_file_id = $3
		WHERE id = $1
	`

	queryGetVoiceForDownloadTranscription = `
		SELECT id, response_file_id
		FROM voice_recognize_req
		WHERE process_step = $1
	`

	querySetStatusDownload = `
		UPDATE voice_recognize_req 
		SET process_step = $2, transcription = $3
		WHERE id = $1
	`

	queryGetVoiceForCreateSummary = `
		SELECT id, transcription
		FROM voice_recognize_req
		WHERE process_step = $1
	`
)