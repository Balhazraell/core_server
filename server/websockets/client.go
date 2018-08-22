package websockets

import (
	"fmt"
	"io"
	"sync"

	"golang.org/x/net/websocket"
)

const channalBufSize = 100

var maxID int

type IncomingMessage struct {
	// По ходу это не входящее сообщение а просто формат передачи сообщений, один и туда и обратно!
	HandlerName string `json:"handler_name"`
	Data        string `json:"data"`
}

type Client struct {
	id int // Должен браться из базы в соответствии с id пользователя в базе.
	ws *websocket.Conn
	ch chan string // Канал для отправки сообщения.

	shutdownRead chan bool
	// shutdownWrite chan bool
}

func (client *Client) Shutdown() {
	client.shutdownRead <- true
	// client.shutdownWrite <- true
}

func (client *Client) Listen() {
	var wg sync.WaitGroup
	// wg.Add(2)
	wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	client.listenWrite()
	// }()

	go func() {
		defer wg.Done()
		client.listenRead()
	}()

	wg.Wait()

	fmt.Printf("Listen у клиента %d работу закончил \n", client.id)
	AppServer.DelClient(client)
}

func (client *Client) SetGameMap(gameMap []byte) {
	// Получили карту - отправляем её пользователю.
	newMessage := IncomingMessage{
		HandlerName: "set_grid",
		Data:        string(gameMap),
	}

	websocket.JSON.Send(client.ws, newMessage)
}

// func (client *Client) listenWrite() {
// 	//
// 	defer func() {
// 		fmt.Printf("listenWrite у клиента %d работу закончил \n", client.id)
// 	}()

// 	for {
// 		select {
// 		case <-client.shutdownWrite:
// 			return
// 		}
// 	}
// }

func (client *Client) listenRead() {
	defer func() {
		fmt.Printf("listenRead у клиента %d работу закончил \n", client.id)
	}()

	for {
		select {
		case <-client.shutdownRead:
			return
		default:
			var msg IncomingMessage
			err := websocket.JSON.Receive(client.ws, &msg)

			if err == io.EOF {
				// client.shutdownWrite <- true
				return
			} else if err != nil {
				fmt.Printf("Проблема чтения сообщения от клиента \n: %v", err)
				// client.shutdownWrite <- true
				return
			} else {
				AppServer.IncomingMessage(client, &msg)
			}
		}
	}
}
