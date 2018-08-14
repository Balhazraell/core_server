package websockets

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

// пока так, но потом надо сделать отдельную инициализацию...
var AppServer Server

type Server struct {
	// TODO: Непонятно зачем нужен pattern в данном случае.
	clients map[int]*Client

	// Каналы
	doneCh chan bool
	// inComing  chan string
	outComing chan string
}

func Start() {
	fmt.Println("Websocket start...")
	clients := make(map[int]*Client)
	doneCh := make(chan bool)
	// inComing := make(chan string)
	outComing := make(chan string)

	AppServer = Server{
		clients,
		doneCh,
		// inComing,
		outComing,
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
		client.Listen()

	}

	http.Handle("/appgame", websocket.Handler(onConnected))
	for {
		select {
		case <-server.doneCh:
			fmt.Println("AAAAAAAAAAAAAA!")
			return
		// case <-server.inComing:
		// 	fmt.Println("Пришло сообщение от пользователя.")
		case <-server.outComing:
			fmt.Println("Необходимо отослать сообщение пользователям.")
		}
	}
	fmt.Println("BBBBBBBBBBBB!")
}

// func (server *Server) onConnected(ws *websocket.Conn) {
// 	defer func() {
// 		err := ws.Close()
// 		if err == nil {
// 			fmt.Println("Websocker close error!")
// 			fmt.Println(err)
// 		}
// 	}()
// 	client := server.newClient(ws)
// 	server.clients[client.id] = client
// 	client.Listen()

// 	fmt.Println("New client is connected")
// }

func (server *Server) newClient(ws *websocket.Conn) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	ch := make(chan string, channalBufSize)
	doneCh := make(chan bool)
	client := &Client{maxID, ws, ch, doneCh}
	maxID++

	return client
}

func (server *Server) DelClient(client *Client) {
	delete(AppServer.clients, client.id)
	fmt.Println("Клиент удален!")
}

func (server *Server) IncomingMessage(client *Client, message *IncomingMessage) {
	fmt.Println("Пришло сообщение от пользователя.")
}
