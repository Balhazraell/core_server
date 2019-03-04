package room

import (
	"encoding/json"
	"fmt"

	"../../logger"

	"github.com/streadway/amqp"
)

var Room RoomStruct

type Chunc struct {
	ID          int      `json:"id"`
	State       int      `json:"state"`
	Сoordinates [][2]int `json:"coordinates"`
}

type Client struct {
	ID int
}

// Room это игровое пространство/"Карта" в котором происходит действие.
type RoomStruct struct {
	ID        int
	Map       map[int]*Chunc
	clients   map[int]*Client
	brockerCh *amqp.Channel

	// Переменные логики.
	GameState int // Делаем крестики нолики, по этому 2 состояния - ходит один потом другой.

	// Каналы
	shutdownLoop chan bool
}

func createMap() {
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

			Room.Map[chunckIDCounter] = &chunc
			chunckIDCounter++
		}
		y += step
	}
}

// StartNewRoom - метод запуска новой комнаты. - В последующем это будет делать main
// На вход подается id комнаты котурую надо создать.
func StartNewRoom(id int) {
	Room := RoomStruct{
		ID:           id,
		Map:          make(map[int]*Chunc),
		clients:      make(map[int]*Client),
		GameState:    1,
		shutdownLoop: make(chan bool),
	}

	createMap()
	StartRabbitMQ(fmt.Sprintf("room_", id))
	go Room.loop()
}

// Stop - Метод принадлежит Room.
// Служит для прекращения работы комнаты.
func (room *RoomStruct) Stop() {
	// Какая-нибудь логика завершения работы.
	room.shutdownLoop <- true
}

func (room *RoomStruct) loop() {
	defer func() {
		logger.InfoPrintf("Комната с id=%v закончила работу.", room.ID)
	}()

	logger.InfoPrintf("Комната с id=%v начала работу.", room.ID)

	for {
		// Обновление логики происходит тут.

		select {
		case <-room.shutdownLoop:
			return
		}
	}
}

func updateAllClientsMap() {
	gameMap, err := json.Marshal(Room.Map)

	if err != nil {
		logger.ErrorPrintf("При формировнии json из игроой карты, произошла ошибка %v.", err)
		return
	}

	clientsIDs := make([]int, 0, len(Room.clients))
	for k := range Room.clients {
		clientsIDs = append(clientsIDs, k)
	}

	newUpdateClientsMapStruct := UpdateClientsMapStruct{
		Map:        gameMap,
		ClientsIDs: clientsIDs,
	}

	updateClientsMapStructJson, err := json.Marshal(newUpdateClientsMapStruct)
	if err != nil {
		logger.ErrorPrintf("При формировнии json из шаблона для отпраки сообщения произошла ошибка %v.", err)
		return
	}

	newMessage := Message{
		HandlerName: "UpdateClientsMap",
		Data:        updateClientsMapStructJson,
	}

	PublishMessage(newMessage)
}

//--------------------- Обработка API -----------------------//
func clientConnect(clientID int) {
	// Необходимо добавить в комнату пользователя.
	logger.InfoPrintf("К комнате %v подключился новый клиент с id=%v.", room.ID, clientID)
	newClient := Client{ID: clientID}
	Room.clients[clientID] = &newClient

	// TODO: Раньше при подключении пользователю отдавалась карта, теперь надо придумать другой способ
	// gameMap, err := json.Marshal(room.Map)
	// return gameMap
}

func clientDisconnect(clientID int) {
	_, ok := Room.clients[clientID]
	if ok {
		logger.InfoPrintf("Удаляем клиента id=%v из комнаты id=%v.", clientID, room.ID)
		delete(Room.clients, clientID)
	} else {
		logger.WarningPrintf("Попытка удалить клиента из комнаты, корого уже нет: id=%v.", clientID)
	}
}

// SetChunckState - Метод принадлежит Room.
// Метод вызываемый при попытке пользователя что-то сделать с участком карты.
func SetChunckState(clientID int, chunkID int) {
	if Room.Map[chunkID].State == 0 {
		Room.Map[chunkID].State = Room.GameState

		if Room.GameState == 1 {
			Room.GameState = 2
		} else {
			Room.GameState = 1
		}

		updateAllClientsMap()
	} else {
		logger.WarningPrintf("Попытка изменить значение в поле с изменненым значеним клиентом с id=%v.", clientID)
		GameServer.SendErrorToСlient(clientID, "Нельзя изменить значение!")
	}
}
