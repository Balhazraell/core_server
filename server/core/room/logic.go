package room

import (
	"encoding/json"

	"../../logger"
)

// Функция создания карты для комнаты.
func createMap() {
	var step = 100
	var y = 0
	var chunckIDCounter = 0

	for i := 0; i < 3; i++ {
		var x = 0
		for j := 0; j < 3; j++ {
			chunc := Chunc{
				ID:    chunckIDCounter,
				State: ChuncStateEmpty,
			}
			chunc.Сoordinates = append(
				chunc.Сoordinates,
				[2]int{x, y},
				[2]int{x + step, y},
				[2]int{x + step, y + step},
				[2]int{x, y + step},
				// необходимо указывать 5 элементов, так как последняя точка замыкает фигуру,
				// это нужно для определения пересечения координат мышки и элементов.
				[2]int{x, y},
			)

			x += step

			Room.Map[chunckIDCounter] = &chunc
			chunckIDCounter++
		}
		y += step
	}
}

func updateClientsMap(clientsIDs []int) {
	logger.InfoPrint("Обновление карт пользователей.")
	gameMap, err := json.Marshal(Room.Map)

	if err != nil {
		logger.WarningPrintf("При формировнии json при упаковке карты, для обновления у клиентов произошла ошибка %v.", err)
		return
	}

	updateMap := UpdateMapStruct{
		Map:        gameMap,
		ClientsIDs: clientsIDs,
	}

	data, err := json.Marshal(updateMap)

	if err != nil {
		logger.WarningPrintf("При формировнии json при создании сооббщения для обновления карт произошла ошибка %v.", err)
		return
	}

	newMessage := MessageRMQ{
		HandlerName: "UpdateClientsMap",
		Data:        string(data),
	}

	PublishMessage(newMessage)
}

//--------------------- Обработка API -----------------------//
func clientConnect(clientID int) {
	logger.InfoPrintf("К комнате %v подключился новый клиент с id=%v.", Room.ID, clientID)

	var elementIndex = findElementInArray(Room.clients, clientID)
	var callbackMessage CallbackMessageStruct

	if elementIndex == -1 {
		Room.clients = append(Room.clients, clientID)
		callbackMessage = CallbackMessageStruct{
			Status:  true,
			Message: "",
		}

		updateClientsMap([]int{clientID})
	} else {
		callbackMessage = CallbackMessageStruct{
			Status:  false,
			Message: "Пользователь с таким id уже есть!",
		}
	}

	message, err := json.Marshal(callbackMessage)
	if err != nil {
		logger.WarningPrintf("При по попытке сформировать json для callback произошла ошибка: %v", err)
		return
	}

	messageRMQ := MessageRMQ{
		HandlerName: "ClientConnectCallback",
		Data:        string(message),
	}

	PublishMessage(messageRMQ)
}

func clientDisconnect(clientID int) {
	var clientIndex = findElementInArray(Room.clients, clientID)

	if clientIndex != -1 {
		logger.InfoPrintf("Удаляем клиента id=%v из комнаты id=%v.", clientID, Room.ID)
		Room.clients = deleElementFromArraByIndex(Room.clients, clientIndex)
	} else {
		logger.WarningPrintf("Попытка удалить клиента из комнаты, корого нет: id=%v.", clientID)
	}
}

// SetChunckState - Метод вызываемый при попытке пользователя что-то сделать с участком карты.
func SetChunckState(clientID int, chunkID int) {
	if Room.Map[chunkID].State == ChuncStateEmpty {
		Room.Map[chunkID].State = Room.GameState

		if Room.GameState == ChuncStateCross {
			Room.GameState = ChuncStateZero
		} else {
			Room.GameState = ChuncStateCross
		}

		updateClientsMap(Room.clients)
	} else {
		logger.WarningPrintf("Попытка изменить значение в поле с изменненым значеним клиентом с id=%v.", clientID)
		//TODO: надо справочник ошибок с кодами ошибок и в коде работать только с кодами ошибок.
		sendErrorMessageStruct := SendErrorMessageStruct{
			ClientID:     clientID,
			ErrorMessage: "Нельзя изменить значение!",
		}

		data, err := json.Marshal(sendErrorMessageStruct)
		if err != nil {
			logger.WarningPrintf("При формировнии json при формировании сообщения с ошибкой произошла ошибка %v.", err)
			return
		}

		newMessage := MessageRMQ{
			HandlerName: "SendErrorMessage",
			Data:        string(data),
		}

		PublishMessage(newMessage)
		//TODO: необходимо отослать сообщение об ошибке поведения пользователя.
		// GameServer.SendErrorToСlient(clientID, "Нельзя изменить значение!")
	}
}

//------------------------------- tools -------------------------//
func deleElementFromArraByIndex(array []int, index int) []int {
	array[index] = array[len(array)-1]
	return array[:len(array)-1]
}

func findElementInArray(array []int, val int) int {
	var result = -1

	for i, v := range array {
		if val == v {
			return i
		}
	}

	return result
}
