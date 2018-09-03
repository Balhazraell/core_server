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
			server.newConnect(clietnID)
		case clientID := <-api.API.ClientDisconnectChl:
			server.clientDisconnect(clientID)
		case chunckStateData := <-api.API.SetChunckStateChl:
			server.setChunckState(
				chunckStateData.ClientID,
				chunckStateData.ChuncID,
			)
		}
	}
}

func (server *gameServer) newConnect(clietnId int) {
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
func (server *gameServer) setChunckState(clientID int, chuncID int) {
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
		fmt.Printf("Удаляем клиента с сервера: id=%v. \n", clientID)
		client.Room.ClientDisconnect(clientID)
		delete(server.Clients, clientID)
	} else {
		fmt.Printf("Попытка удалить клиента корого уже нет: id=%v. \n", clientID)
	}
}
