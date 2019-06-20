package core

import (
	"encoding/json"
	"fmt"

	"github.com/Balhazraell/logger"
)

// APIMetods - Перечень доступных API методов.
var APIMetods = map[string]func(string, int){
	"RoomConnect":      apiRoomConnect,
	"UpdateClientsMap": apiUpdateClientsMap,
	"SendErrorMessage": apiSendErrorMessage,
}

//------------------ Income struct ---------------------------//
type updateMapStruct struct {
	Map        []byte `json:"Map"`
	ClientsIDs []int  `json:"ClientsIDs"`
}

type sendErrorMessageStruct struct {
	ClientID     int    `json:"ClientID"`
	ErrorMessage string `json:"ErrorMessage"`
}

//--------------------- Outgoing struct -------------------------//
type setChunckStateStruct struct {
	ClientID int `json:"ClientID"`
	ChunkID  int `json:"ChunkID"`
}

//--------------------- Обработка API -----------------------//
func apiRoomConnect(data string, roomID int) {
	var ID int
	err := json.Unmarshal([]byte(data), &ID)

	if err != nil {
		logger.ErrorPrintf("Ошибка распаковки JSON: \nОшибка: %v \nДанные: %v", err, data)
	}

	status, message := validateRoomConnect(ID)
	callbackMessage := callbackStruct{
		RoomID:  -1,
		UserID:  -1,
		Status:  status,
		Message: message,
	}

	CreateMessage(fmt.Sprintf("room_%v", ID), callbackMessage, "CallbackRoomConnect")

	if status {
		newRoomConnect(ID)
	}
}

func apiUpdateClientsMap(data string, roomID int) {
	var object updateMapStruct
	err := json.Unmarshal([]byte(data), &object)

	if err != nil {
		logger.ErrorPrintf("Ошибка распаковки JSON: \nОшибка: %v \nДанные: %v", err, data)
	}

	var success_clients []int

	// Вырожденный случай!
	for i := 0; i < len(object.ClientsIDs); i++ {
		status, message := validateUpdateClientsMap(object.ClientsIDs[i])
		if status {
			success_clients = append(success_clients, object.ClientsIDs[i])
		} else {
			callbackMessage := callbackStruct{
				RoomID:  -1,
				UserID:  object.ClientsIDs[i],
				Status:  status,
				Message: message,
			}

			// Отправляем по одному проблемному клиенту.
			CreateMessage(fmt.Sprintf("room_%v", roomID), callbackMessage, "CallbackUpdateClientsMap")
		}
	}

	if len(success_clients) > 0 {
		updateClientsMap(object.Map, object.ClientsIDs)
	}
}

func apiSendErrorMessage(data string, roomID int) {
	var object sendErrorMessageStruct
	err := json.Unmarshal([]byte(data), &object)

	if err != nil {
		logger.ErrorPrintf("Ошибка распаковки JSON: \nОшибка: %v \nДанные: %v", err, data)
	}

	status, message := validateSendErrorMessage(object.ClientID)
	callbackMessage := callbackStruct{
		RoomID:  -1,
		UserID:  object.ClientID,
		Status:  status,
		Message: message,
	}

	CreateMessage(fmt.Sprintf("room_%v", roomID), callbackMessage, "CallbackSendErrorMessage")

	if status {
		sendErrorMessage(object.ClientID, object.ErrorMessage)
	}
}
