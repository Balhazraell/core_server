package core

import (
	"../api"
	"../logger"
)

// В первой итерации вебсокеты будут передавать сообщения на прямую в gameServer.
// но потом надо продумать другую связь, возможно через балансировщик.

// GameServer - singletone для работы с care частью сервера.GameServer
var GameServer gameServer

// Client структура описывающая связь клиента и комнаты в которой он находится.
// Считается что пользователь не может быть в не комнат - тоесть хотя бы в какой-то он точно есть.
type Client struct {
	ID   int
	Room *Room
}

type gameServer struct {
	Clients map[int]*Client
	Rooms   map[int]*Room

	shutdownLoop chan bool
}

// GameServerStart - метод запуска игрового сервера.
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

// Stop - метод завершения работы игрового сервера.
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

func (server *gameServer) newConnect(clietnID int) {
	// сейчас пока буду закидывать в первую комнату.
	// Подключаем по id комнаты в которую он входит.
	// TODO: ЭТОНАДО РАЗДЕЛИТЬ НА ДВЕ ФУНКЦИИ!!!!
	// подключение нового пользователя и получение им карты это два разных события.

	logger.InfoPrint("На сервер пришел новый пользователь.")

	room, ok := server.Rooms[1]
	if !ok {
		logger.WarningPrintf("Попытка присоединится к комнате которой нет: Клиет - %v", clietnID)
		return
	}

	client := Client{
		clietnID,
		room,
	}

	GameServer.Clients[client.ID] = &client

	newClientIsConnectedStruct := api.NewClientIsConnectedStruct{
		ClientID:     clietnID,
		ClientMap:    room.ClientConnect(&client),
		RoomsCatalog: server.getRoomsIDsList(),
	}

	api.API.NewClientIsConnectedChl <- newClientIsConnectedStruct
}

func (server *gameServer) setChunckState(clientID int, chuncID int) {
	logger.InfoPrint("На сервер пришело сообщение об обновлении состояния комнаты")
	server.Clients[clientID].Room.SetChunckState(clientID, chuncID)
}

func (server *gameServer) UpdateClientsMap(gameMap []byte, clientsIDs []int) {
	updateClientsMapStruct := api.UpdateClientsMapStruct{
		GameMap:    gameMap,
		ClientsIDs: clientsIDs,
	}

	api.API.UpdateClientsMapChl <- updateClientsMapStruct
}

func (server *gameServer) SendErrorToСlient(clientID int, message string) {
	sendErrorToСlientStruct := api.SendErrorToСlientStruct{
		ClientID: clientID,
		Message:  message,
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

		var clientsIDs = []int{clietn.ID}

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
