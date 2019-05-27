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
		message = fmt.Sprintf("API метода %v в комнате room_id:%v не существует.", handlerName, Room.ID)
	}

	return status, message
}
