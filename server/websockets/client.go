package websockets

import (
	"encoding/json"
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

type OutcomingMessage struct {
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
	wg.Add(1)

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
	fmt.Printf("Пытаемся задать карту клиенту с id = %v.\n ", client.id)

	newMessage := OutcomingMessage{
		HandlerName: "set_grid",
		Data:        string(gameMap),
	}

	fmt.Printf("newMessage сформирован и готовится к отправке.\n")

	websocket.JSON.Send(client.ws, newMessage)
}

func (client *Client) listenRead() {
	defer func() {
		fmt.Printf("listenRead у клиента %d работу закончил \n", client.id)
	}()

	// TODO: убрать это отсюда.
	// api.API.ClientConnectionChl <- ClientMaxId

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
				fmt.Printf("Проблема чтения сообщения от клиента : %v \n", err)
				// client.shutdownWrite <- true
				return
			} else {
				AppServer.IncomingMessage(client, &msg)
			}
		}
	}
}

func (client *Client) SendError(message string) {
	jsonMessage, err := json.Marshal(message)

	if err != nil {
		fmt.Printf("При формировнии json при отправки сообщения об ошибке произошла ошибка %v \n", err)
	}

	newMessage := OutcomingMessage{
		HandlerName: "send_error",
		Data:        string(jsonMessage),
	}

	websocket.JSON.Send(client.ws, newMessage)
}

func (client *Client) SetRoomsCatalog(roomsIDs []int) {
	jsonRoomsIDs, err := json.Marshal(roomsIDs)

	if err != nil {
		fmt.Printf("Произошла ошибка при формироваии json в SetRoomsCatalog: %v \n", err)
	}

	newMessage := OutcomingMessage{
		HandlerName: "set_rooms_catalog",
		Data:        string(jsonRoomsIDs),
	}

	websocket.JSON.Send(client.ws, newMessage)
}
