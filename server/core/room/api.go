package room

import (
	"encoding/json"

	"../../logger"
)

type CallbackMessageStruct struct {
	ServiceID int    `json:"ServiceID"`
	Status    bool   `json:"Status"`
	Message   string `json:"Message"`
}

//--------------------- core struct --------------------//
type UpdateMapStruct struct {
	Map        []byte `json:"Map"`
	ClientsIDs []int  `json:"ClientsIDs"`
}

type SendErrorMessageStruct struct {
	ClientID     int    `json:"ClientID"`
	ErrorMessage string `json:"ErrorMessage"`
}

//--------------------- room struct --------------------//
type SetChunckStateStruct struct {
	ClientID int `json:"ClientID"`
	ChunkID  int `json:"ChunkID"`
}

// APIMetods - Перечень доступных API методов.
var APIMetods = map[string]func(string){
	"ClientConnect":    apiClientConnect,
	"ClientDisconnect": apiClientDisconnect,
	"SetChunckState":   apiSetChunckState,
}

func apiClientConnect(data string) {
	var clientID int
	err := json.Unmarshal([]byte(data), &clientID)

	if err != nil {
		logger.ErrorPrintf("Ошибка API при подключении нового клиента: %s;\n Ошибка в данных: %s", err, data)
	}

	clientConnect(clientID)
}

func apiClientDisconnect(data string) {
	var clientID int
	err := json.Unmarshal([]byte(data), &clientID)

	if err != nil {
		logger.ErrorPrintf("Ошибка API при отключении клиента: %s;\n Ошибка в данных: %s", err, data)
	}

	clientDisconnect(clientID)
}

func apiSetChunckState(data string) {
	var setChunckStateStruct SetChunckStateStruct

	err := json.Unmarshal([]byte(data), &setChunckStateStruct)

	if err != nil {
		logger.ErrorPrintf("Ошибка API задании состояния чанку: %s;\n Ошибка в данных: %s", err, data)
	}

	SetChunckState(setChunckStateStruct.ClientID, setChunckStateStruct.ChunkID)
}
