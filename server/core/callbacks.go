package core

import (
	"encoding/json"

	"github.com/Balhazraell/logger"
)

type callbackStruct struct {
	RoomID  int    `json:"RoomID"`
	UserID  int    `json:"UserID"`
	Status  bool   `json:"Status"`
	Message string `json:"Message"`
}

// CallbackMetods - Перечень методов для получения ответов при запросах.
var CallbackMetods = map[string]func(string){
	"СallbackAPICall":          сallbackAPICall,
	"CallbackClientConnect":    сallbackСlientConnect,
	"CallbackClientDisconnect": сallbackСlientDisconnect,
}

func сallbackAPICall(data string) {
	var callback = callbackStruct{}
	err := json.Unmarshal([]byte(data), &callback)

	if err != nil {
		logger.ErrorPrintf("Ошибка распаковки JSON: \nОшибка: %v \nДанные: %v", err, data)
	}

	if !callback.Status {
		logger.ErrorPrintf("Ошибка вызова API метода: \n%v", callback.Message)
	}
}

func сallbackСlientConnect(data string) {
	var object clientConnectCallbackStruct
	err := json.Unmarshal([]byte(data), &object)

	if err != nil {
		logger.ErrorPrintf("Ошибка при распаковке данных callback при подключении клиента: %s", err)
	}

	// clientConnectCallback(object.ClientID, object.Status, object.Message)
}

func сallbackСlientDisconnect(data string) {
	var object clientDisconnectCallbackStruct
	err := json.Unmarshal([]byte(data), &object)

	if err != nil {
		logger.ErrorPrintf("Ошибка при распаковке данных callback при подключении клиента: %s", err)
	}

	// clientDisconectCallback(object.ClientID, object.Status, object.Message)
}
