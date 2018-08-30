package core

import (
	"fmt"

	"../api"
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

	// создадим пока одну комнату, для тестирования.
	room1 := StartNewRoom(666)
	GameServer.Rooms[room1.ID] = room1
}

func Stop() {
	GameServer.shutdownLoop <- true
}

func (server *gameServer) loop() {
	defer func() {
		fmt.Println("Игровой сервер закончил свою работу.")
	}()

	fmt.Println("Игровой сервер запущен.")

	for {
		// Что-нибудь что делают всякие сервера.

		select {
		case <-server.shutdownLoop:
			return
		case clietnID := <-api.API.ClientConnectionChl:
			server.NewConnect(clietnID)
		case clientID := <-api.API.ClientDisconnectChl:
			fmt.Println("Клиент %v отключился", clientID)
		case chunckStateData := <-api.API.SetChunckStateChl:
			server.SetChunckState(
				chunckStateData.ClientID,
				chunckStateData.ChuncID,
			)
		}
	}
}

func (server *gameServer) NewConnect(clietnId int) {
	// сейчас пока буду закидывать в первую комнату.
	// Подключаем по id комнаты в которую он входит.
	// TODO: ЭТОНАДО РАЗДЕЛИТЬ НА ДВЕ ФУНКЦИИ!!!!
	// подключение нового пользователя и получение им карты это два разных события.

	fmt.Println("На сервер пришел новый пользователь")

	room := GameServer.Rooms[666]

	client := Client{
		clietnId,
		room,
	}

	GameServer.Clients[client.Id] = &client

	newClientIsConnectedStruct := api.NewClientIsConnectedStruct{
		clietnId,
		room.ClientConnect(&client),
	}

	api.API.NewClientIsConnectedChl <- newClientIsConnectedStruct
}

// Интерфейсы для получения данных от
func (server *gameServer) SetChunckState(clientID int, chuncID int) {
	fmt.Println("На сервер пришело сообщение об обновлении состояния комнаты")
	server.Clients[clientID].Room.SetChunckState(clientID, chuncID)
}

func (server *gameServer) UpdateClientsMap(gameMap []byte, clientsIDs []int) {
	updateClientsMapStruct := api.UpdateClientsMapStruct{
		gameMap,
		clientsIDs,
	}

	api.API.UpdateClientsMapChl <- updateClientsMapStruct
}
