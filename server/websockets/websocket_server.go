package websockets

import (
	"fmt"
	"net/http"
	"runtime"

	"golang.org/x/net/websocket"

	"../core"
)

// пока так, но потом надо сделать отдельную инициализацию...
var AppServer Server

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

	CoreMetods map[string]func(int, []byte)
}

func Start() {
	fmt.Println("Websocket start...")
	clients := make(map[int]*Client)
	shutdownCh := make(chan bool)
	// inComing := make(chan string)
	outComing := make(chan string)

	CoreMetods := map[string]func(int, []byte){
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

		client := server.newClient(ws)
		server.clients[client.id] = client
		fmt.Println("New client is connected")
		fmt.Printf("%v goroutine is running \n", runtime.NumGoroutine())
		client.Listen()

	}

	http.Handle("/appgame", websocket.Handler(onConnected))
	for {
		select {
		case <-server.shutdownCh:
			return
		case <-server.outComing:
			fmt.Println("Необходимо отослать сообщение пользователям.")
		}
	}
}

func (server *Server) newClient(ws *websocket.Conn) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	// сделать через две разные функции...
	// Блять... вообще отдачу карты должна инициировать комната...
	clientId, gameMap := core.GameServer.NewConnect(666)
	ch := make(chan string, channalBufSize)
	shutdownRead := make(chan bool)
	// shutdownWrite := make(chan bool)
	client := &Client{clientId, ws, ch, shutdownRead}

	client.SetGameMap(gameMap)

	return client
}

func (server *Server) DelClient(client *Client) {
	delete(AppServer.clients, client.id)
	fmt.Printf("Client whith id %v is deleted \n", client.id)
	fmt.Printf("%v goroutine is running \n", runtime.NumGoroutine())
}

func (server *Server) IncomingMessage(client *Client, message *IncomingMessage) {
	// скорее всего надо не сразу дергать методы игрового сервера, а нормально распарсить их тут и
	// вызывать конкретные методы с конкретными аргументами.
	server.CoreMetods[message.HandlerName](client.id, message.Data)
}

func setChunckState(clientID int, data []byte) {
	// var chunckStateStructure ChunckStateStructure

	// json.Unmarshal(data, &chunckStateStructure)

	fmt.Println(data)

	// chunckID, ok := data[0].(int)

	// if !ok {
	// 	fmt.Println("Не смогли распарсить данные setChunckState")
	// }
	// core.SetChunckState(clientID, chunckID)
}
