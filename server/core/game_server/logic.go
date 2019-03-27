package game_server

import (
	"fmt"

	"../../api"
	"../../logger"
)

func newRoomConnect(id int) {
	logger.InfoPrintf("Подключается новая комната с id:%v", id)
	_, ok := Server.Rooms[id]
	if !ok {
		Server.Rooms[id] = fmt.Sprintf("room_%v", id)
	} else {
		logger.ErrorPrintf("Комната с id=%v уже существует!", id)
	}
}

//--------------------- Обработка сообщений от клиента -----------------------//
func clientConnect(clietnID int) {
	logger.InfoPrint("На сервер пришел новый пользователь.")

	// TODO:
	/*
		Должен братся набор существующих комнат и смотреть в какую мы можем поместить пользователя...
		Если такой комнаты нет, то нужно создать новую.
	*/
	_, ok := Server.RoomIDByClient[clietnID]
	if !ok {
		// TODO: Комнаты может ещё/уже не быть!
		Server.RoomIDByClient[clietnID] = 1
		logger.InfoPrint(Server.Rooms)
		CreateMessage(Server.Rooms[1], clietnID, "ClientConnect")
	} else {
		logger.WarningPrintf("Пользователь с id:%v уже существует.", clietnID)
	}
}

func setChunckState(clientID int, chuncID int) {
	logger.InfoPrint("На сервер пришело сообщение об обновлении состояния комнаты")

	message := setChunckStateStruct{
		ClientID: clientID,
		ChunkID:  chuncID,
	}

	roomID := Server.RoomIDByClient[clientID]

	CreateMessage(Server.Rooms[roomID], message, "SetChunckState")
}

func changeRoom(clientID int, newRoomID int) {
	_, ok := Server.Rooms[newRoomID]
	if ok {
		// var currentRoomID = Server.RoomIDByClient[clientID]

		//! Тут отправляется сообщение о попытке отключить пользователя от комнаты.

		//! Сообщение о подключении нового пользователя.

		// clietn.Room = room

		// var clientsIDs = []int{clietn.ID}

		// server.UpdateClientsMap(
		// 	room.ClientConnect(clietn),
		// 	clientsIDs,
		// )
	} else {
		logger.WarningPrintf("Err: Клиент %v пытается подключится к комнате которой нет: %v.", clientID, newRoomID)
		// TODO: нужно отправить сообщение клиенту, о том, что произошла ошибка.
	}
}

func clientDisconnect(clientID int) {
	_, ok := Server.RoomIDByClient[clientID]
	if ok {
		logger.InfoPrintf("Удаляем клиента с сервера: id=%v.", clientID)
		//! Отправляем сообщение о попытке удалить пользователя из комнаты.

		delete(Server.RoomIDByClient, clientID)
	} else {
		logger.WarningPrintf("Попытка удалить клиента корого уже нет: id=%v.", clientID)
		return
	}
}

// getRoomsIDAndName - Возвращает список имен комнат
func getRoomsData() []api.RoomData {
	var roomsData []api.RoomData

	for k, v := range Server.Rooms {
		newRoomData := api.RoomData{
			ID:   k,
			Name: v,
		}

		roomsData = append(roomsData, newRoomData)
	}

	return roomsData
}

//--------------------- Обработка API -----------------------//
func updateClientsMap(gameMap []byte, clientsIDs []int) {
	updateClientsMapStruct := api.UpdateClientsMapStruct{
		GameMap:    gameMap,
		ClientsIDs: clientsIDs,
	}

	api.API.UpdateClientsMapChl <- updateClientsMapStruct
}

func sendErrorMessage(clientID int, message string) {
	sendErrorToСlientStruct := api.SendErrorToСlientStruct{
		ClientID: clientID,
		Message:  message,
	}

	api.API.SendErrorToСlientChl <- sendErrorToСlientStruct
}

func clientConnectCallback(clientID int, status bool, message string) {
	if status {
		newClientIsConnectedStruct := api.NewClientIsConnectedStruct{
			ClientID:     clientID,
			RoomsCatalog: getRoomsData(),
		}

		api.API.NewClientIsConnectedChl <- newClientIsConnectedStruct

	} else {
		delete(Server.RoomIDByClient, clientID)
	}
}

func clientDisconectCallback(clientID int, status bool, message string) {
	if status {
		delete(Server.RoomIDByClient, clientID)
	}
}
