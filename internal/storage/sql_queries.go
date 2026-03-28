package storage

const (
	querySaveIncommingVoice = `
		INSERT INTO voice_recognize_req (voice_data, process_step) 
        VALUES ($1, 'BEGIN')
		RETURNING id
	`
)