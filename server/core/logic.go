package core

import (
	"fmt"

	"../api"
	"github.com/Balhazraell/logger"
	"github.com/Balhazraell/tools"
)

func newRoomConnect(id int) {
	logger.InfoPrintf("Подключается новая комната с id:%v", id)
	_, ok := Server.Rooms[id]
	if !ok {
		Server.Rooms[id] = fmt.Sprintf("room_%v", id)

		//!TODO сейчас всех кидаем в первую появившуюся комнату.
		if len(Server.PendingUsers) != 0 {
			for len(Server.PendingUsers) != 0 {
				clientID := Server.PendingUsers[0]
				Server.PendingUsers = tools.DeleElementFromArraByIndex(Server.PendingUsers, 0)

				logger.InfoPrintf("Пользователя с id:%v из очереди ожидания поместили в комнату id:%v", clientID, id)

				Server.RoomIDByClient[clientID] = id
				CreateMessage(Server.Rooms[id], clientID, "ClientConnect")
			}
		} else {
			updateRoomsCatalog()
		}
	} else {
		logger.ErrorPrintf("Комната с id=%v уже существует!", id)
	}
}

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

func getAllUsers() []int {
	result := make([]int, 0)
	for k := range Server.RoomIDByClient {
		result = append(result, k)
	}

	return result
}

func updateRoomsCatalog() {
	roomsCatalog := api.RoomsCatalogStruct{
		ClientIDs:    getAllUsers(),
		RoomsCatalog: getRoomsData(),
	}

	api.API.UpdateRoomsCatalog <- roomsCatalog
}

//--------------------- Обработка сообщений от клиента -----------------------//
func clientConnect(clientID int) {
	logger.InfoPrint("На сервер пришел новый пользователь.")

	// TODO:
	/*
		Должен братся набор существующих комнат и смотреть в какую мы можем поместить пользователя...
		Если такой комнаты нет, то нужно создать новую.
		!!! Нет обработки отключения комнаты.
	*/
	_, ok := Server.RoomIDByClient[clientID]
	if !ok {
		//TODO: Комнаты может ещё/уже не быть!
		keys := make([]int, 0)

		for key := range Server.Rooms {
			keys = append(keys, key)
		}

		if len(keys) == 0 {
			logger.InfoPrint("Еще не создано ни одной комнаты, пользователь находится в ожидании!")
			//!TODO: Необходимо отправить пользователю окно ожидания.
			Server.RoomIDByClient[clientID] = -1
			Server.PendingUsers = append(Server.PendingUsers, clientID)
			return
		}

		//Server.RoomIDByClient[clientID] = keys[0]
		CreateMessage(Server.Rooms[keys[0]], clientID, "ClientConnect")
	} else {
		logger.WarningPrintf("Пользователь с id:%v уже существует.", clientID)
	}
}

func setChunckState(clientID int, chuncID int) {
	logger.InfoPrint("На сервер пришело сообщение об обновлении состояния комнаты")

	message := setChunckStateStruct{
		ClientID: clientID,
		ChunkID:  chuncID,
	}

	roomID := Server.RoomIDByClient[clientID]
	logger.InfoPrint(Server.Rooms[roomID])
	CreateMessage(Server.Rooms[roomID], message, "SetChunckState")
}

func changeRoom(clientID int, newRoomID int) {
	_, ok := Server.Rooms[newRoomID]
	if ok {
		logger.InfoPrintf("Новая комната: %v", newRoomID)

		oldRoomID := Server.RoomIDByClient[clientID]

		CreateMessage(Server.Rooms[oldRoomID], clientID, "ClientDisconnect")
		CreateMessage(Server.Rooms[newRoomID], clientID, "ClientConnect")

		Server.RoomIDByClient[clientID] = newRoomID
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

func clientConnectCallback(roomID int, clientID int, status bool, message string) {
	if status {
		//?TODO: отправка списка комнат идет отдельно, возможно этот вызо не нужен.
		newClientIsConnectedStruct := api.NewClientIsConnectedStruct{
			ClientID:     clientID,
			RoomsCatalog: getRoomsData(),
		}

		/*
			Есть очеред ожидающих подсоединения!
			Нужно человека добавлять туда и только после ответа добавлять в список на сервере с комнатами.
			Причем если пришел ответ что мы не можем подключится - пытаемся бросить в другую комнату.
			Такая же тема с переключающимися. В сообщении от комнаты должна быть вся информация.

			Бредовая мысль но может не надо отдавать карту и список комнат и что бы то ни было ещё, пока их не попросят...
			тоесть - мы сообщаем что Пользователь подключился комнате, а дальше внешние системы пусть запрашивают что им надо...

			! Для сервисной архитектуры... мы кладем сообщение о новом пользователе в список сообщений и карта отправится по запросу "ударасердца"
		*/

		api.API.NewClientIsConnectedChl <- newClientIsConnectedStruct

	} else {
		delete(Server.RoomIDByClient, clientID)
		Server.PendingUsers = append(Server.PendingUsers, clientID)
	}
}

func clientDisconectCallback(roomID int, clientID int, status bool, message string) {
	if status {
		delete(Server.RoomIDByClient, clientID)
	}
}
