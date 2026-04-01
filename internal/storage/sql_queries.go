package storage

const (
	querySaveIncomingVoice = `
		INSERT INTO voice_recognize_req (voice_data, chat_id, process_step) 
        VALUES ($1, $2, $3)
		RETURNING id
	`

	queryGetVoicesForUpload = `
		SELECT id, chat_id, voice_data
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
		SELECT id, chat_id, request_file_id
		FROM voice_recognize_req
		WHERE process_step = $1
	`

	querySetStatusRecognition = `
		UPDATE voice_recognize_req 
		SET process_step = $2, recognize_id = $3
		WHERE id = $1
	`

	queryGetVoiceForCheckStatus = `
		SELECT id, chat_id, recognize_id
		FROM voice_recognize_req
		WHERE process_step = $1
	`

	querySetStatusWait = `
		UPDATE voice_recognize_req 
		SET process_step = $2, response_file_id = $3
		WHERE id = $1
	`

	queryGetVoiceForDownloadTranscription = `
		SELECT id, chat_id, response_file_id
		FROM voice_recognize_req
		WHERE process_step = $1
	`

	querySetStatusDownload = `
		UPDATE voice_recognize_req 
		SET process_step = $2, transcription = $3
		WHERE id = $1
	`

	queryGetVoiceForCreateSummary = `
		SELECT id, chat_id, transcription
		FROM voice_recognize_req
		WHERE process_step = $1
	`

	querySetStatusSuccess = `
		UPDATE voice_recognize_req
		SET process_step = $2, summary = $3, created_at = $4
		WHERE id = $1
	`

	queryGetVoiceByID = `
		SELECT summary
		FROM voice_recognize_req
		WHERE id = $1 AND chat_id = $2
	`

	queryGetListSummaryID = `
		SELECT id
		FROM voice_recognize_req
		WHERE chat_id = $1
		ORDER BY created_at
	`
)