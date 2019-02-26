package core

import (
	"encoding/json"
	"fmt"
	"io"

	"../api"
	"../logger"

	"github.com/streadway/amqp"
)

// В первой итерации вебсокеты будут передавать сообщения на прямую в gameServer.
// но потом надо продумать другую связь, возможно через балансировщик.

// GameServer - singletone для работы с care частью сервера.GameServer
var GameServer gameServer

// CoreMetods - набор методов вызываемых у core части.(по сути это надор API)
var CoreMetods = map[string]func(string){
	"roomConnect": GameServer.RoomConnect,
}

// Client структура описывающая связь клиента и комнаты в которой он находится.
// Считается что пользователь не может быть в не комнат - тоесть хотя бы в какой-то он точно есть.
type Client struct {
	ID     int
	RoomID int
}

type RoomIdAndNameStruct struct {
	ID   int
	Name string
}

type gameServer struct {
	Clients map[int]*Client
	// Rooms          map[int]*Room
	QeueRoomsNames map[int]string

	shutdownLoop chan bool
}

// TODO: потом надо перемеименовать в просто message.
type coreMessage struct {
	HandlerName string `json:"handler_name"`
	Data        string `json:"data"`
}

// GameServerStart - метод запуска игрового сервера.
func GameServerStart() {
	var clients = make(map[int]*Client)
	// var rooms = make(map[int]*Room)
	var qeueRoomsNames = make(map[int]string)
	var shutdownLoop = make(chan bool)

	GameServer = gameServer{
		clients,
		// rooms,
		qeueRoomsNames,
		shutdownLoop,
	}

	go GameServer.loop()

	//--------------------------------- Owerall ----------------------------
	// Создадим связь с брокером.
	conn, err := amqp.Dial("amqp://macroserv:12345@localhost:5672/macroserv")
	if err != nil {
		logger.ErrorPrintf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Устанавливаем соединение с брокером.
	ch, err := conn.Channel()
	if err != nil {
		logger.ErrorPrintf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	// Точка доступа должна быть создана, до того как создана очередь.
	// так как слать сообщения в несучествующую точку доступа запрещено!
	err = ch.ExchangeDeclare(
		"core",   // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		logger.ErrorPrintf("Failed to declare an exchange: %s", err)
	}

	err = ch.ExchangeDeclare(
		"rooms",  // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		logger.ErrorPrintf("Failed to declare an exchange: %s", err)
	}

	//--------------------------------- For core ----------------------------
	// Создаем очередь из которой будем поулчать сообщения.
	// Делается всегда и там где принимается и там где отправляется,
	// если очереди нет то сообщение просто проигнорится,
	// но если очередь оздана хотя бы раз, то повторно создана не будет.
	// так как это очередь для того что бы слушать сообщения приходящие нам,
	// не надо его запоминать, у нас будет горутина крутится...
	queue, err := ch.QueueDeclare(
		"сore", // name
		true,   // durable
		false,  // delete when usused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		logger.ErrorPrintf("Failed to declare a queue: %s", err)
	}

	err = ch.QueueBind(
		queue.Name, // queue name
		queue.Name, // routing key (binding_key)
		// TODO: наверно надо вынести в отдельную переменную.
		"core", // exchange
		false,
		nil,
	)
	if err != nil {
		logger.ErrorPrintf("Failed to bind a queue: %s", err)
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

	// Запускаем горутину которая будет "слушать" очередь.
	go func() {
		for d := range msgs {
			var msg coreMessage
			err := json.Unmarshal(d.Body, &msg)

			if err == io.EOF {
				continue
			} else if err != nil {
				logger.ErrorPrintf("Проблема чтения сообщения от комнаты : %v.", err)
				continue
			} else {
				// TODO: можем упасть при вызове не верного метода - надо обработать!
				// Допустим метод которого нет в списке.
				// TODO: Написать тесты для этого метода.
				CoreMetods[msg.HandlerName](msg.Data)
			}
		}
	}()

	// TODO: Сейчас создаем несколько комнат.
	// В сервисной архитектуре, это будем делать не мы.
	StartNewRoom(1)
	StartNewRoom(2)
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

	// room, ok := server.Rooms[1]
	// if !ok {
	// 	logger.WarningPrintf("Попытка присоединится к комнате которой нет: Клиет - %v", clietnID)
	// 	return
	// }

	// TODO:
	// Тут должен быть метод выдающий id свободной комнаты.
	// Необходимо подумать можно ли динамически растить новые комнаты.
	// Можно создать определенное количество комнат и добавлять пользователей как зрителей.

	client := Client{
		clietnID,
		1,
	}

	GameServer.Clients[client.ID] = &client

	// TODO: мы должны послать запрос комнате можно ли в неё подключится и если нет, то попытаться подключтья к другой комнате.
	// Если свободных комнат нет, надо будет их создавать...
	newClientIsConnectedStruct := api.NewClientIsConnectedStruct{
		ClientID:     clietnID,
		ClientMap:    room.ClientConnect(&client),
		RoomsCatalog: server.getRoomsIDAndName(),
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

// getRoomsIDAndName - Возвращает список имен комнат
func (server *gameServer) getRoomsIDAndName() []RoomIdAndNameStruct {
	// TODO: должна возвращаться структура [[id, name]]
	roomsIDAndNames := make([]RoomIdAndNameStruct, 0, len(server.QeueRoomsNames))
	for k := range server.QeueRoomsNames {
		newRoomIDAndNameStruct := RoomIdAndNameStruct{ID: k, Name: server.QeueRoomsNames[k]}
		roomsIDAndNames = append(roomsIDAndNames, newRoomIDAndNameStruct)
	}

	return roomsIDAndNames
}

// --------------------------------------------------------
// Тут будет набор методов вызываемый через RabbitMQ потом это должно стать API этого модуля.
func (server *gameServer) RoomConnect(message string) {
	// TODO: должно передоваться имя и ID комнаты... (Надо подумать над этим...)
	var ID int
	err := json.Unmarshal([]byte(message), &ID)

	if err != nil {
		logger.ErrorPrintf("Failed unmarshal room connect message: %s", err)
	}
	// TODO не понятно на сколько хорошая практика.
	GameServer.QeueRoomsNames[ID] = fmt.Sprintf("room_%d", ID)
}
