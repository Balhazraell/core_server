package core

import (
	"fmt"
)

type Chunc struct {
	State       int
	Сoordinates [][2]int
}

type Room struct {
	Id      int
	Map     []*Chunc
	clients map[int]*Client

	shutdownLoop chan bool
}

func (room *Room) createMap() {
	var step = 100
	var y = 0

	for i := 0; i < 3; i++ {
		var x = 0
		for j := 0; j < 3; j++ {
			chunc := Chunc{}
			chunc.Сoordinates = append(
				chunc.Сoordinates,
				[2]int{x, y},
				[2]int{x + step, y},
				[2]int{x + step, y + step},
				[2]int{x, y + step},
			)

			x += step

			room.Map = append(room.Map, &chunc)
		}
	}
}

func StartNewRoom(id int) *Room {
	newRoom := Room{}
	newRoom.Id = id
	newRoom.clients = make(map[int]*Client)
	newRoom.shutdownLoop = make(chan bool)
	newRoom.createMap()

	go newRoom.loop()

	return &newRoom
}

func (room *Room) Stop() {
	// Какая-нибудь логика завершения работы.

	room.shutdownLoop <- true
}

func (room *Room) ClientConnect(client *Client) []*Chunc {
	// Необходимо добавить в комнату пользователя.
	fmt.Printf("К комнате %v подключился новый клиент с id=%v.", room.Id, client.Id)
	room.clients[client.Id] = client
	return room.Map
}

func (room *Room) ClientDisconnect(client *Client) {
	fmt.Printf("Из комнта %v вышел клиент с ")
	delete(room.clients, client.Id)
}

func (room *Room) loop() {
	defer func() {
		fmt.Printf("Комната с id=v% закончила работу.", room.Id)
	}()

	fmt.Printf("Комната с id=v% начала работу.", room.Id)

	for {
		// Обновление логики происходит тут.

		select {
		case <-room.shutdownLoop:
			return
		}
	}
}
