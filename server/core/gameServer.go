package core

import (
	"../api"
	"../logger"
)

// В первой итерации вебсокеты будут передавать сообщения на прямую в gameServer.
// но потом надо продумать другую связь, возможно через балансировщик.

var GameServer gameServer

type Client struct {
	Id   int
	Room *Room
}

type gameServer struct {
	Clients map[int]*Client
	Rooms   map[int]*Room

	shutdownLoop chan bool
}

func GameServerStart() {
	var clients = make(map[int]*Client)
	var rooms = make(map[int]*Room)
	var shutdownLoop = make(chan bool)

	GameServer = gameServer{
		clients,
		rooms,
		shutdownLoop,
	}

	go GameServer.loop()

	// создадим несколько комнат, для тестирования.
	room1 := StartNewRoom(1)
	room2 := StartNewRoom(2)
	GameServer.Rooms[room1.ID] = room1
	GameServer.Rooms[room2.ID] = room2

}

func Stop() {
	GameServer.shutdownLoop <- true
}

func (server *gameServer) loop() {
	defer func() {
		logger.InfoPrint("Игровой сервер закончил свою работу.")
	}()

	logger.InfoPrint("Игровой сервер запущен.")

	for {
		select {
		case <-server.shutdownLoop:
			return
		case clietnID := <-api.API.ClientConnectionChl:
			server.newConnect(clietnID)
		case clientID := <-api.API.ClientDisconnectChl:
			server.clientDisconnect(clientID)
		case chunckStateData := <-api.API.SetChunckStateChl:
			server.setChunckState(
				chunckStateData.ClientID,
				chunckStateData.ChuncID,
			)
		case changeRoomStructureData := <-api.API.ChangeRoomChl:
			server.changeRoom(
				changeRoomStructureData.ClientID,
				changeRoomStructureData.RoomID,
			)
		}
	}
}

func (server *gameServer) newConnect(clietnId int) {
	// сейчас пока буду закидывать в первую комнату.
	// Подключаем по id комнаты в которую он входит.
	// TODO: ЭТОНАДО РАЗДЕЛИТЬ НА ДВЕ ФУНКЦИИ!!!!
	// подключение нового пользователя и получение им карты это два разных события.

	logger.InfoPrint("На сервер пришел новый пользователь.")

	room, ok := server.Rooms[1]
	if !ok {
		logger.WarningPrintf("Попытка присоединится к комнате которой нет: Клиет - %v", clietnId)
		return
	}

	client := Client{
		clietnId,
		room,
	}

	GameServer.Clients[client.Id] = &client

	newClientIsConnectedStruct := api.NewClientIsConnectedStruct{
		clietnId,
		room.ClientConnect(&client),
		server.getRoomsIDsList(),
	}

	api.API.NewClientIsConnectedChl <- newClientIsConnectedStruct
}

// Интерфейсы для получения данных от
func (server *gameServer) setChunckState(clientID int, chuncID int) {
	logger.InfoPrint("На сервер пришело сообщение об обновлении состояния комнаты")
	server.Clients[clientID].Room.SetChunckState(clientID, chuncID)
}

func (server *gameServer) UpdateClientsMap(gameMap []byte, clientsIDs []int) {
	updateClientsMapStruct := api.UpdateClientsMapStruct{
		gameMap,
		clientsIDs,
	}

	api.API.UpdateClientsMapChl <- updateClientsMapStruct
}

func (server *gameServer) SendErrorToСlient(client_id int, message string) {
	sendErrorToСlientStruct := api.SendErrorToСlientStruct{
		client_id,
		message,
	}

	api.API.SendErrorToСlientChl <- sendErrorToСlientStruct
}

func (server *gameServer) clientDisconnect(clientID int) {
	client, ok := server.Clients[clientID]
	if ok {
		logger.InfoPrintf("Удаляем клиента с сервера: id=%v.", clientID)
		client.Room.ClientDisconnect(clientID)
		delete(server.Clients, clientID)
	} else {
		logger.WarningPrintf("Попытка удалить клиента корого уже нет: id=%v.", clientID)
		return
	}
}

func (server *gameServer) changeRoom(clientID int, newRoomID int) {
	var clietn = server.Clients[clientID]

	clietn.Room.ClientDisconnect(clientID)

	room, ok := server.Rooms[newRoomID]
	if ok {
		clietn.Room = room

		var clientsIDs = []int{clietn.Id}

		server.UpdateClientsMap(
			room.ClientConnect(clietn),
			clientsIDs,
		)
	} else {
		logger.WarningPrintf("Err: Клиент %v пытается подключится к комнате которой нет: %v.", clientID, newRoomID)
		// TODO: нужно отправить сообщение клиенту, о том, что произошла ошибка.
		delete(server.Clients, clientID)
	}
}

func (server *gameServer) getRoomsIDsList() []int {
	roomsIDs := make([]int, 0, len(server.Rooms))
	for k := range server.Rooms {
		roomsIDs = append(roomsIDs, k)
	}

	return roomsIDs
}
