package models

const (
	MsgInternalServerError = "Ошибка на сервере. Повторите запрос позже."
	MsgFailSummaryProcess = "Не удалось обработать встречу. ID: %d. Повторите отправку."
	MsgSuccesSummaryProcess = "Встреча успешно сохранена! Вы можете получить ее по ID: %d."
	MsgErrSaveVoice = "Не удалось загрузить встречу. Повторите отправку."
	MsgSuccessSaveVoice = "Встреча принята в обработку. ID: %d."
	MsgErrInvalidFormat = "Пожалуйста, отправьте аудиофайл в формате MP3."

	GetErrEmptyID = "Пожалуйста, укажите ID встречи.\nПример: /get 123"
	GetErrInvalidID = "Неверный формат ID. ID должен быть числом."
	Get404 = "По вашему запросу ничего не найдено."
)