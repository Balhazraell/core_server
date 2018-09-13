package websockets

import (
	"encoding/json"
	"net/http"
	"runtime"

	"golang.org/x/net/websocket"

	"../api"
	"../logger"
)

// пока так, но потом надо сделать отдельную инициализацию...
var AppServer Server

// Когда клиент подсоединяется - должен передоваться его id...
var ClientMaxId int

type ChunckStateStructure struct {
	ChunckID int `json:"chunck_id"`
}

type ChangeRoomStructure struct {
	RoomId int `json:"room_id"`
}

type Server struct {
	// TODO: Непонятно зачем нужен pattern в данном случае.
	clients map[int]*Client

	// Каналы
	shutdownCh chan bool

	CoreMetods map[string]func(int, string)
}

func Start() {
	logger.InfoPrint("Websocket start...")
	clients := make(map[int]*Client)
	shutdownCh := make(chan bool)

	CoreMetods := map[string]func(int, string){
		"setChunckState": setChunckState,
		"chengeRoomID":   chengeRoomID,
	}

	AppServer = Server{
		clients,
		shutdownCh,
		CoreMetods,
	}
	go AppServer.listen()
}

func (server *Server) listen() {
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				logger.ErrorPrintf("Websocker закрылся с ошибкой: %v", err)
			}
			logger.InfoPrint("Websocker закрыт...")
		}()

		server.newClient(ws)
	}

	http.Handle("/appgame", websocket.Handler(onConnected))
	for {
		select {
		case <-server.shutdownCh:
			return
		case updateClientsMapStruct := <-api.API.UpdateClientsMapChl:
			server.updateClientsMap(
				updateClientsMapStruct.GameMap,
				updateClientsMapStruct.ClientsIDs,
			)
		case connectedClientData := <-api.API.NewClientIsConnectedChl:
			server.setClietMap(
				connectedClientData.ClientID,
				connectedClientData.ClientMap,
			)
			server.setRoomsCatalog(
				connectedClientData.ClientID,
				connectedClientData.RoomsCatalog,
			)
		case sendErrorToСlientStruct := <-api.API.SendErrorToСlientChl:
			server.sendErrorToСlient(
				sendErrorToСlientStruct.ClientID,
				sendErrorToСlientStruct.Message,
			)
		}
	}
}

func (server *Server) newClient(ws *websocket.Conn) {
	if ws == nil {
		logger.ErrorPrint("ws не может быть равен nil!")
		return
	}

	ch := make(chan string, channalBufSize)
	shutdownRead := make(chan bool)
	client := &Client{ClientMaxId, ws, ch, shutdownRead}

	server.clients[client.id] = client

	logger.InfoPrint("Новый клиент присоединился.")
	// TODO: на основе этого принта надо сделать тест на корректную диструктуризацию горутин.
	logger.InfoPrintf("%v горутин сейчас запущенно.", runtime.NumGoroutine())

	// Создали канал, запустили его, теперь можно и игровому серверпусказать что подключился игрок.
	api.API.ClientConnectionChl <- ClientMaxId
	ClientMaxId++
	// надо наверно сделать так что бы вызов этого метода не тормазил работу метода
	client.Listen()
}

func (server *Server) DelClient(client *Client) {
	api.API.ClientDisconnectChl <- client.id
	delete(AppServer.clients, client.id)
	logger.InfoPrintf("Клиент с id %v удален.", client.id)
	// TODO: на основе этого принта надо сделать тест на корректную диструктуризацию горутин.
	logger.InfoPrintf("%v горутин сейчас запущенно", runtime.NumGoroutine())
}

func (server *Server) IncomingMessage(client *Client, message *IncomingMessage) {
	// скорее всего надо не сразу дергать методы игрового сервера, а нормально распарсить их тут и
	// вызывать конкретные методы с конкретными аргументами.

	server.CoreMetods[message.HandlerName](client.id, message.Data)
}

func setChunckState(clientID int, data string) {
	var chunckStateStructure ChunckStateStructure
	err := json.Unmarshal([]byte(data), &chunckStateStructure)

	if err != nil {
		logger.WarningPrintf("Ошибка парсинга json в setChunckState %v.", err)
		return
	}

	setChunckStateStruct := api.SetChunckStateStruct{
		clientID,
		chunckStateStructure.ChunckID,
	}

	api.API.SetChunckStateChl <- setChunckStateStruct
}

func (server *Server) updateClientsMap(gameMap []byte, clientsIDs []int) {
	// TODO: необходимы тесты на то корректность удаление данных и отработку даже с некорретными пришедшими данными.
	// Типа если пришел id которого нет в списке id-шников.
	for i := 0; i < len(clientsIDs); i++ {
		server.setClietMap(clientsIDs[i], gameMap)
	}
}

func (server *Server) setClietMap(clietID int, clientMap []byte) {
	client, ok := server.clients[clietID]
	if ok {
		client.SetGameMap(clientMap)
	} else {
		logger.WarningPrintf("Попытка задать карту клиенту которого уже нет: %v.", clietID)
		return
	}
}

func (server *Server) setRoomsCatalog(clietID int, roomsIDs []int) {
	client, ok := server.clients[clietID]
	if ok {
		client.SetRoomsCatalog(roomsIDs)
	} else {
		logger.WarningPrintf("Попытка задать карту клиенту которого уже нет: %v.", clietID)
		return
	}
}

func (server *Server) sendErrorToСlient(clientID int, message string) {
	server.clients[clientID].SendError(message)
}

func chengeRoomID(clientID int, data string) {
	var changeRoomStructure ChangeRoomStructure

	err := json.Unmarshal([]byte(data), &changeRoomStructure)

	if err != nil {
		logger.WarningPrintf("Ошибка парсинга json в setChunckState %v", err)
	}

	changeRoomStructureForCore := api.ChangeRoomStructure{
		clientID,
		changeRoomStructure.RoomId}

	api.API.ChangeRoomChl <- changeRoomStructureForCore
}
