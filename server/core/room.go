package core

import (
	"encoding/json"
	"fmt"
	"log"

	"../logger"

	"github.com/streadway/amqp"
)

/*
	У комнат несколько методов которые по сути должны быть API:
	StartNewRoom,
	Stop,
	ClientConnect,
	ClientDisconnect,
	SetChunckState,
*/

// Chunc описывает струтуру участка игрового пространства.
type Chunc struct {
	ID          int      `json:"id"`
	State       int      `json:"state"`
	Сoordinates [][2]int `json:"coordinates"`
}

// Room это игровое пространство/"Карта" в котором происходит действие.
// Комната живет своей жизнью.
// Комната состоит из частей (Chunc)
type Room struct {
	ID        int
	Map       map[int]*Chunc
	clients   map[int]*Client
	brockerCh *amqp.Channel

	// Переменные логики.
	GameState int // Делаем крестики нолики, по этому 2 состояния - ходит один потом другой.

	// Каналы
	shutdownLoop chan bool
	updateMap    chan bool
}

func (room *Room) createMap() {
	var step = 100
	var y = 0
	var chunckIDCounter = 0

	for i := 0; i < 3; i++ {
		var x = 0
		for j := 0; j < 3; j++ {
			chunc := Chunc{}
			chunc.ID = chunckIDCounter
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

			room.Map[chunckIDCounter] = &chunc
			chunckIDCounter++
		}
		y += step
	}
}

// StartNewRoom - метод запуска новой комнаты.
// На вход подается id комнаты котурую надо создать.
func StartNewRoom(id int) *Room {
	newRoom := Room{}
	newRoom.ID = id
	newRoom.Map = make(map[int]*Chunc)
	newRoom.clients = make(map[int]*Client)
	newRoom.GameState = 1
	newRoom.shutdownLoop = make(chan bool)
	newRoom.createMap()

	go newRoom.loop()

	// Сейчас создадим полноценное соединение для RabbitMQ
	conn, err := amqp.Dial("amqp://macroserv:12345@localhost:15672/")
	if err != nil {
		logger.ErrorPrintf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Создаем канал.
	ch, err := conn.Channel()
	if err != nil {
		logger.ErrorPrintf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	// Сначала создаем очередь на получение сообщений, назвние
	// будет формироваться из имени комнаты, в нашем случае из id
	queue, err := ch.QueueDeclare(
		fmt.Sprintf("room_%d", id), // name
		false,                      // durable
		false,                      // delete when usused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		logger.ErrorPrintf("Failed to declare a queue: %s", err)
	}

	// Теперь создаем подписчика.
	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		logger.ErrorPrintf("Failed to register a consumer: %s", err)
	}

	// Мониторим очередь на наличие сообщений.
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	return &newRoom
}

// Stop - Метод принадлежит Room.
// Служит для прекращения работы комнаты.
func (room *Room) Stop() {
	// Какая-нибудь логика завершения работы.
	room.shutdownLoop <- true
}

// ClientConnect - Метод принадлежит Room.
// Метод подключения нового пользователя к комнате.
func (room *Room) ClientConnect(client *Client) []byte {
	// Необходимо добавить в комнату пользователя.
	logger.InfoPrintf("К комнате %v подключился новый клиент с id=%v.", room.ID, client.ID)
	room.clients[client.ID] = client
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

// SetChunckState - Метод принадлежит Room.
// Метод вызываемый при попытке пользователя что-то сделать с участком карты.
func (room *Room) SetChunckState(clientID int, chunkID int) {
	if room.Map[chunkID].State == 0 {
		room.Map[chunkID].State = room.GameState

		if room.GameState == 1 {
			room.GameState = 2
		} else {
			room.GameState = 1
		}

		room.updateClientsMap()
	} else {
		logger.WarningPrintf("Попытка изменить значение в поле с изменненым значеним клиентом с id=%v.", clientID)
		GameServer.SendErrorToСlient(clientID, "Нельзя изменить значение!")
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

// ClientDisconnect метод отключающий пользователя от этой комнаты.
func (room *Room) ClientDisconnect(clientID int) {
	_, ok := room.clients[clientID]
	if ok {
		logger.InfoPrintf("Удаляем клиента id=%v из комнаты id=%v.", clientID, room.ID)
		delete(room.clients, clientID)
	} else {
		logger.WarningPrintf("Попытка удалить клиента из комнаты, корого уже нет: id=%v.", clientID)
	}
}
