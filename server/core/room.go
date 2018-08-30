package core

import (
	"encoding/json"
	"fmt"
)

type Chunc struct {
	ID          int      `json:"id"`
	State       int      `json:"state"`
	Сoordinates [][2]int `json:"coordinates"`
}

type Room struct {
	ID      int
	Map     map[int]*Chunc
	clients map[int]*Client

	// Переменные логики.
	GameState int // Делаем крестики нолики, по этому 2 состояния - ходит один потом другой.

	// Каналы
	shutdownLoop chan bool
	updateMap    chan bool
}

func (room *Room) createMap() {
	var step = 100
	var y = 0
	var chunckIdCounter = 0

	for i := 0; i < 3; i++ {
		var x = 0
		for j := 0; j < 3; j++ {
			chunc := Chunc{}
			chunc.ID = chunckIdCounter
			chunc.Сoordinates = append(
				chunc.Сoordinates,
				[2]int{x, y},
				[2]int{x + step, y},
				[2]int{x + step, y + step},
				[2]int{x, y + step},
				// необходимо указывать 5 элементов, так как последняя точка замыкает фигуру,
				// это нужно для определения пересечения координат мышки и элементов.
				[2]int{x, y},
			)

			x += step

			room.Map[chunckIdCounter] = &chunc
			chunckIdCounter++
		}
		y += step
	}
}

func StartNewRoom(id int) *Room {
	newRoom := Room{}
	newRoom.ID = id
	newRoom.Map = make(map[int]*Chunc)
	newRoom.clients = make(map[int]*Client)
	newRoom.GameState = 1
	newRoom.shutdownLoop = make(chan bool)
	newRoom.createMap()

	go newRoom.loop()

	return &newRoom
}

func (room *Room) Stop() {
	// Какая-нибудь логика завершения работы.
	room.shutdownLoop <- true
}

func (room *Room) ClientConnect(client *Client) []byte {
	// Необходимо добавить в комнату пользователя.
	fmt.Printf("К комнате %v подключился новый клиент с id=%v.\n", room.ID, client.Id)
	room.clients[client.Id] = client
	// ВОобще не при подключении надо возвращать карту, это надо делать по специальной функции,
	// Наверно надо отдавать в loop в канал id пользователя, кому надо задать карту...
	gameMap, err := json.Marshal(room.Map)

	if err != nil {
		fmt.Printf("При формировнии json при подключении нового клиента произошла ошибка %v \n", err)
	}

	return gameMap
}

func (room *Room) ClientDisconnect(client *Client) {
	fmt.Printf("Из комнта %v вышел клиент с id = v% \n", room.ID, client.Id)
	delete(room.clients, client.Id)
}

func (room *Room) loop() {
	defer func() {
		fmt.Printf("Комната с id=v% закончила работу. \n", room.ID)
	}()

	fmt.Printf("Комната с id=v% начала работу. \n", room.ID)

	for {
		// Обновление логики происходит тут.

		select {
		case <-room.shutdownLoop:
			return

		// Даже не знаю на сколько целесообразно делать это в отдельном потоке.
		// Мсль была в том, что update карт должен произоти не моментально после изменений
		// но хз на сколько это грамотоное решение.
		case <-room.updateMap:
			room.updateClientsMap()
		}

	}
}

func (room *Room) SetChunckState(client_id int, chunk_id int) {
	fmt.Println("Задано состояние для чанка.")
	room.Map[chunk_id].State = room.GameState
	if room.GameState == 1 {
		room.GameState = 2
	} else {
		room.GameState = 1
	}

	room.updateClientsMap()
}

func (room *Room) updateClientsMap() {
	fmt.Println("Обновление карт пользователей.")
	gameMap, err := json.Marshal(room.Map)

	if err != nil {
		fmt.Printf("При формировнии json при подключении нового клиента произошла ошибка %v \n", err)
	}

	clientsIDs := make([]int, 0, len(room.clients))
	for k := range room.clients {
		clientsIDs = append(clientsIDs, k)
	}

	GameServer.UpdateClientsMap(gameMap, clientsIDs)
}
