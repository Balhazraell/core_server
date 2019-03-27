package room

import (
	"fmt"

	"../../logger"
)

// Room - основная структура данных этого пакета.
var Room room

type chunc struct {
	ID          int      `json:"id"`
	State       int      `json:"state"`
	Сoordinates [][2]int `json:"coordinates"`
}

type room struct {
	ID      int
	Map     map[int]*chunc
	clients []int

	// Переменные логики.
	GameState int // Делаем крестики нолики, по этому 2 состояния - ходит один потом другой.

	// Каналы
	shutdownLoop chan bool
	updateMap    chan bool
}

// StartNewRoom - метод запуска новой комнаты.
// На вход подается id комнаты котурую надо создать.
func StartNewRoom(id int) {
	Room := room{
		ID:           id,
		Map:          make(map[int]*chunc),
		shutdownLoop: make(chan bool),
		updateMap:    make(chan bool),
	}

	createMap()

	StartRabbitMQ(fmt.Sprintf("room_%v", id))
	go Room.loop()
}

// Stop - Останавлием работу комнаты
func (r *room) Stop() {
	// ...какая-нибудь логика завершения работы.
	r.shutdownLoop <- true
}

func (r *room) loop() {
	defer func() {
		logger.InfoPrintf("Комната с id=%v закончила работу.", r.ID)
	}()

	logger.InfoPrintf("Комната с id=%v начала работу.", r.ID)

	for {
		// Обновление логики происходит тут.

		select {
		case <-r.shutdownLoop:
			return

		// Даже не знаю на сколько целесообразно делать это в отдельном потоке.
		// Мсль была в том, что update карт должен произоти не моментально после изменений
		// но хз на сколько это грамотоное решение.
		case <-r.updateMap:
			updateClientsMap(Room.clients)
		}
	}
}
