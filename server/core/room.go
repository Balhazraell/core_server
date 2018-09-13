package core

import (
	"encoding/json"

	"../logger"
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
	logger.InfoPrintf("К комнате %v подключился новый клиент с id=%v.", room.ID, client.Id)
	room.clients[client.Id] = client
	// ВОобще не при подключении надо возвращать карту, это надо делать по специальной функции,
	// Наверно надо отдавать в loop в канал id пользователя, кому надо задать карту...
	gameMap, err := json.Marshal(room.Map)

	if err != nil {
		// TODO: тут надо сделать так, что бы функция завершилась ничего не возвращая...
		logger.WarningPrintf("При формировнии json при подключении нового клиента произошла ошибка %v.", err)
	}

	return gameMap
}

func (room *Room) loop() {
	defer func() {
		logger.InfoPrintf("Комната с id=%v закончила работу.", room.ID)
	}()

	logger.InfoPrintf("Комната с id=%v начала работу.", room.ID)

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
	if room.Map[chunk_id].State == 0 {
		room.Map[chunk_id].State = room.GameState

		if room.GameState == 1 {
			room.GameState = 2
		} else {
			room.GameState = 1
		}

		room.updateClientsMap()
	} else {
		logger.WarningPrintf("Попытка изменить значение в поле с изменненым значеним клиентом с id=%v.", client_id)
		GameServer.SendErrorToСlient(client_id, "Нельзя изменить значение!")
	}
}

func (room *Room) updateClientsMap() {
	logger.InfoPrint("Обновление карт пользователей.")
	gameMap, err := json.Marshal(room.Map)

	if err != nil {
		logger.WarningPrintf("При формировнии json при подключении нового клиента произошла ошибка %v.", err)
		return
	}

	clientsIDs := make([]int, 0, len(room.clients))
	for k := range room.clients {
		clientsIDs = append(clientsIDs, k)
	}

	GameServer.UpdateClientsMap(gameMap, clientsIDs)
}

func (room *Room) ClientDisconnect(clientID int) {
	_, ok := room.clients[clientID]
	if ok {
		logger.InfoPrintf("Удаляем клиента id=%v из комнаты id=%v.", clientID, room.ID)
		delete(room.clients, clientID)
	} else {
		logger.WarningPrintf("Попытка удалить клиента из комнаты, корого уже нет: id=%v.", clientID)
	}
}
