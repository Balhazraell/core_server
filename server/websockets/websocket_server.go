package websockets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"golang.org/x/net/websocket"

	"../api"
)

// пока так, но потом надо сделать отдельную инициализацию...
var AppServer Server

// Когда клиент подсоединяется - должен передоваться его id...
var ClientMaxId int

type ChunckStateStructure struct {
	ChunckID int `json:"chunck_id"`
}

type Server struct {
	// TODO: Непонятно зачем нужен pattern в данном случае.
	clients map[int]*Client

	// Каналы
	shutdownCh chan bool
	// inComing  chan string
	outComing chan string

	CoreMetods map[string]func(int, string)
}

func Start() {
	fmt.Println("Websocket start...")
	clients := make(map[int]*Client)
	shutdownCh := make(chan bool)
	// inComing := make(chan string)
	outComing := make(chan string)

	CoreMetods := map[string]func(int, string){
		"setChunckState": setChunckState,
	}

	AppServer = Server{
		clients,
		shutdownCh,
		// inComing,
		outComing,
		CoreMetods,
	}
	go AppServer.listen()
}

func (server *Server) listen() {
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				fmt.Println("Websocker close error!")
				fmt.Println(err)
			}
			fmt.Println("Websocker close...")
		}()

		server.newClient(ws)
	}

	http.Handle("/appgame", websocket.Handler(onConnected))
	for {
		select {
		case <-server.shutdownCh:
			return
		case <-server.outComing:
			fmt.Println("Необходимо отослать сообщение пользователям.")
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
		panic("ws cannot be nil")
	}

	ch := make(chan string, channalBufSize)
	shutdownRead := make(chan bool)
	client := &Client{ClientMaxId, ws, ch, shutdownRead}

	server.clients[client.id] = client
	fmt.Println("New client is connected")
	fmt.Printf("%v goroutine is running \n", runtime.NumGoroutine())

	// Создали канал, запустили его, теперь можно и игровому серверпусказать что подключился игрок.
	api.API.ClientConnectionChl <- ClientMaxId
	ClientMaxId++
	// надо наверно сделать так что бы вызов этого метода не тормазил работу метода
	client.Listen()
}

func (server *Server) DelClient(client *Client) {
	api.API.ClientDisconnectChl <- client.id
	delete(AppServer.clients, client.id)
	fmt.Printf("Client whith id %v is deleted \n", client.id)
	fmt.Printf("%v goroutine is running \n", runtime.NumGoroutine())
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
		fmt.Println("Ошибка парсинга json в setChunckState %v", err)
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
		fmt.Printf("Err: Попытка задать карту клиенту которого уже нет: %v. \n", clietID)
	}
}

func (server *Server) sendErrorToСlient(clientID int, message string) {
	server.clients[clientID].SendError(message)
}
