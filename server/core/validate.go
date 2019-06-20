package core

import (
	"fmt"

	"github.com/Balhazraell/logger"
)

func validateAPIcall(handlerName string) (bool, string) {
	status := true
	message := ""

	_, ok := Server.APIandCallbackMetods[handlerName]
	if !ok {
		logger.WarningPrintf("Попытка вызвать API которго нет или к которому нет доступа: %v.", handlerName)
		status = false
		message = fmt.Sprintf("API метода %v в Core нет или нет прав на его использование.", handlerName)
	}

	return status, message
}

func validateRoomConnect(roomID int) (bool, string) {
	status := true
	message := ""

	// Проверка идентификатора на существование.
	_, ok := Server.Rooms[roomID]
	if ok {
		logger.WarningPrintf("Попытка подключения комнаты с уже существующим идентификатором %v", roomID)
		status = false
		message = "Комната с идентификатором уже существует!"
		return status, message
	}

	return status, message
}

func validateUpdateClientsMap(userID int) (bool, string) {
	status := true
	message := ""

	// Проверка идентификатора на существование.
	_, ok := Server.RoomIDByClient[userID]
	if !ok {
		logger.WarningPrintf("Попытка обновить карту у пользователя которого нет среди пользователей %v", userID)
		status = false
		message = "Пользователя с таким id нет!"
	}

	return status, message
}

func validateSendErrorMessage(userID int) (bool, string) {
	status := true
	message := ""

	_, ok := Server.RoomIDByClient[userID]
	if !ok {
		logger.WarningPrintf("Попытка отправить ошибку пользователю которого нет среди пользователей %v", userID)
		status = false
		message = "Пользователя с таким id нет!"
	}

	return status, message
}
