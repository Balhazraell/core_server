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
	// сейчас пока буду закидывать в первую комнату.
	// Подключаем по id комнаты в которую он входит.

	logger.InfoPrint("На сервер пришел новый пользователь.")

	// room, ok := server.Rooms[1]
	// if !ok {
	// 	logger.WarningPrintf("Попытка присоединится к комнате которой нет: Клиет - %v", clietnID)
	// 	return
	// }

	// TODO:
	// Тут должен быть метод выдающий id свободной комнаты.
	// Необходимо подумать можно ли динамически растить новые комнаты.
	// Можно создать определенное количество комнат и добавлять пользователей как зрителей.

	Server.RoomIDByClient[clietnID] = 1

	// TODO: мы должны послать запрос комнате можно ли в неё подключится и если нет, то попытаться подключтья к другой комнате.
	// Если свободных комнат нет, надо будет их создавать...

	//! Тут отправляем сообщение о попытке подключить клиента.

	//! Данный код должен работать как ответ от сервиса о том, что он подключил клиента.
	// newClientIsConnectedStruct := api.NewClientIsConnectedStruct{
	// 	ClientID:     clietnID,
	// 	ClientMap:    room.ClientConnect(&client),
	// 	RoomsCatalog: getRoomsIDAndName(),
	// }

	// api.API.NewClientIsConnectedChl <- newClientIsConnectedStruct
}

func setChunckState(clientID int, chuncID int) {
	logger.InfoPrint("На сервер пришело сообщение об обновлении состояния комнаты")

	//! Тут отправляем сообщение нужной комнате.
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
func sendRoomsIDAndName() {
	// TODO: должна возвращаться структура [[id, name]]
	// roomsIDAndNames := make([]RoomIdAndNameStruct, 0, len(server.QeueRoomsNames))
	// for k := range server.QeueRoomsNames {
	// 	newRoomIDAndNameStruct := RoomIdAndNameStruct{ID: k, Name: server.QeueRoomsNames[k]}
	// 	roomsIDAndNames = append(roomsIDAndNames, newRoomIDAndNameStruct)
	// }

	// return roomsIDAndNames
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
