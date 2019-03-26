package game_server

import (
	"../../api"
	"../../logger"

	// TODO: временное решение для проерки работы комнат.
	"../room"
)

// В первой итерации вебсокеты будут передавать сообщения на прямую в gameServer.
// но потом надо продумать другую связь, возможно через балансировщик.

// Server - singletone для работы с care частью сервера.GameServer
var Server server

type server struct {
	RoomIDByClient map[int]int
	Rooms          map[int]string

	shutdownLoop chan bool
}

// GServerStart - метод запуска игрового сервера.
func ServerStart() {
	Server = server{
		RoomIDByClient: make(map[int]int),
		Rooms:          make(map[int]string),
		shutdownLoop:   make(chan bool),
	}

	go Server.loop()

	// TODO: Сейчас создаем несколько комнат.
	// В сервисной архитектуре, это будем делать не мы.
	room.StartNewRoom(1)
	room.StartNewRoom(2)
}

// Stop - метод завершения работы игрового сервера.
func (serv *server) Stop() {
	serv.shutdownLoop <- true
}

func (serv *server) loop() {
	defer func() {
		logger.InfoPrint("Игровой сервер закончил свою работу.")
	}()

	logger.InfoPrint("Игровой сервер запущен.")

	for {
		select {
		case <-serv.shutdownLoop:
			return
		case clietnID := <-api.API.ClientConnectionChl:
			clientConnect(clietnID)
		case clientID := <-api.API.ClientDisconnectChl:
			clientDisconnect(clientID)
		case chunckStateData := <-api.API.SetChunckStateChl:
			setChunckState(
				chunckStateData.ClientID,
				chunckStateData.ChuncID,
			)
		case changeRoomStructureData := <-api.API.ChangeRoomChl:
			changeRoom(
				changeRoomStructureData.ClientID,
				changeRoomStructureData.RoomID,
			)
		}
	}
}
