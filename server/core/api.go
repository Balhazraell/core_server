package core

import (
	"encoding/json"

	"github.com/Balhazraell/logger"
)

// APIMetods - Перечень доступных API методов.
var APIMetods = map[string]func(string){
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

type clientConnectCallbackStruct struct {
	ClientID int    `json:"ClientID"`
	Status   bool   `json:"Status"`
	Message  string `json:"Message"`
}

type clientDisconnectCallbackStruct struct {
	ClientID int    `json:"ClientID"`
	Status   bool   `json:"Status"`
	Message  string `json:"Message"`
}

//--------------------- Outgoing struct -------------------------//
type setChunckStateStruct struct {
	ClientID int `json:"ClientID"`
	ChunkID  int `json:"ChunkID"`
}

//--------------------- Обработка API -----------------------//
func apiRoomConnect(data string) {
	var ID int
	err := json.Unmarshal([]byte(data), &ID)

	if err != nil {
		logger.ErrorPrintf("Failed unmarshal room connect message: %s", err)
	}

	newRoomConnect(ID)
}

func apiUpdateClientsMap(data string) {
	var object updateMapStruct
	err := json.Unmarshal([]byte(data), &object)

	if err != nil {
		logger.ErrorPrintf("Ошибка при распаковке данных для обновления карты: %s", err)
	}

	updateClientsMap(object.Map, object.ClientsIDs)
}

func apiSendErrorMessage(data string) {
	var object sendErrorMessageStruct
	err := json.Unmarshal([]byte(data), &object)

	if err != nil {
		logger.ErrorPrintf("Ошибка при распаковке данных отправки сообщения об ошибке: %s", err)
	}

	sendErrorMessage(object.ClientID, object.ErrorMessage)
}
