package game_server

import (
	"encoding/json"

	"../../logger"
)

// Перечень доступных API методов.
var APIMetods = map[string]func(string){
	"RoomConnect":              APIRoomConnect,
	"UpdateClientsMap":         APIUpdateClientsMap,
	"SendErrorMessage":         APISendErrorMessage,
	"ClientConnectCallback":    APIClientConnectCallback,
	"ClientDisconnectCallback": APIClientDisconnectCallback,
}

type UpdateMapStruct struct {
	Map        []byte `json:"Map"`
	ClientsIDs []int  `json:"ClientsIDs"`
}

type SendErrorMessageStruct struct {
	ClientID     int    `json:"ClientID"`
	ErrorMessage string `json:"ErrorMessage"`
}

type CallbackMessageStruct struct {
	Status  bool   `json:"Status"`
	Message string `json:"Message"`
}

//--------------------- Обработка API -----------------------//
func APIRoomConnect(data string) {
	// TODO: должно передоваться имя и ID комнаты... (Надо подумать над этим...)
	var ID int
	err := json.Unmarshal([]byte(data), &ID)

	if err != nil {
		logger.ErrorPrintf("Failed unmarshal room connect message: %s", err)
	}
	// TODO не понятно на сколько хорошая практика.
	newRoomConnect(ID)
}

func APIUpdateClientsMap(data string) {
	var updateMapStruct UpdateMapStruct
	err := json.Unmarshal([]byte(data), &updateMapStruct)

	if err != nil {
		logger.ErrorPrintf("Ошибка при распаковке данных для обновления карты: %s", err)
	}
	// TODO не понятно на сколько хорошая практика.
	updateClientsMap(updateMapStruct.Map, updateMapStruct.ClientsIDs)
}

func APISendErrorMessage(data string) {
	var sendErrorMessageStruct SendErrorMessageStruct
	err := json.Unmarshal([]byte(data), &sendErrorMessageStruct)

	if err != nil {
		logger.ErrorPrintf("Ошибка при распаковке данных отправки сообщения об ошибке: %s", err)
	}
	// TODO не понятно на сколько хорошая практика.
	sendErrorMessage(sendErrorMessageStruct.ClientID, sendErrorMessageStruct.ErrorMessage)
}

func APIClientConnectCallback(data string) {

}

func APIClientDisconnectCallback(data string) {

}
