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
		INSERT INTO voice_recognize_req (voice_data, process_step) 
        VALUES ($1, $2)
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
)