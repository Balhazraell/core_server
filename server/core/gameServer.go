package core

import (
	"fmt"
)

// В первой итерации вебсокеты будут передавать сообщения на прямую в gameServer.
// но потом надо продумать другую связь, возможно через балансировщик.

var GameServer gameServer

// Временное решение отвечающее за увеличение id пользователей подключившихся к приложению.
var clientMaxId int

type Client struct {
	Id     int
	RoomId map[int]*Room
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
	GameServer.Rooms[room1.Id] = room1
}

func Stop() {
	GameServer.shutdownLoop <- true
}

func (server *gameServer) loop() {
	defer func() {
		fmt.Printf("Игровой сервер закончил свою работу.")
	}()

	fmt.Printf("Игровой сервер запущен.")

	for {
		// Что-нибудь что делают всякие сервера.

		select {
		case <-server.shutdownLoop:
			return
		}
	}
}

func (server *gameServer) NewConnect(roomId int) (int, []*Chunc) {
	// сейчас пока буду закидывать в первую комнату.
	// Подключаем по id комнаты в которую он входит.

	client := Client{
		clientMaxId,
		make(map[int]*Room),
	}

	client.RoomId[roomId] = GameServer.Rooms[roomId]
	GameServer.Clients[client.Id] = &client

	clientMaxId++

	// так как пользователь подключился впервый раз, сразу отдадим ему сетку.
	gameMap := GameServer.Rooms[roomId].ClientConnect(&client)

	return client.Id, gameMap
}
