package websockets

import (
	"fmt"
	"io"
	"sync"

	"golang.org/x/net/websocket"
)

const channalBufSize = 100

var maxID int

// type RoutingMessage struct {
// 	HandlerName string `json:"handler_name"`
// 	Data struct `json:"data"`
// }
type IncomingMessage struct {
	HandlerName string `json:"handler_name"`
	Data        string `json:"data"`
}

type Client struct {
	id     int // Должен браться из базы в соответствии с id пользователя в базе.
	ws     *websocket.Conn
	ch     chan string // Канал для отправки сообщения.
	doneCh chan bool   // Канал для завершения работы соединения.

	// FIX: надеюсь временное решение.
	doneRead  chan bool
	doneWrite chan bool
}

func (client *Client) Done() {
	// Нужно убить данный объект и его связи...

	// так как запущены бесконечные циклы для прослушивания каналов,
	// теперь их надо останоить.
	client.doneRead <- true
	client.doneWrite <- true
}

func (client *Client) Listen() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		client.listenWrite()
	}()

	go func() {
		defer wg.Done()
		client.listenRead()
	}()

	wg.Wait()

	fmt.Printf("Listen у клиента %d работу закончил \n", client.id)
	AppServer.DelClient(client)
}

func (client *Client) listenWrite() {
	for {
		select {
		case <-client.doneWrite:
			fmt.Printf("listenWrite у клиента %d работу закончил \n", client.id)
			return
		}
	}
}

func (client *Client) listenRead() {
	defer func() {
		fmt.Printf("listenRead у клиента %d работу закончил \n", client.id)
	}()
	for {
		select {
		case <-client.doneRead:
			return
		default:
			var msg IncomingMessage
			err := websocket.JSON.Receive(client.ws, &msg)

			if err == io.EOF {
				client.doneWrite <- true
				return
			} else if err != nil {
				fmt.Printf("Проблема чтения сообщения от клиента \n: %v", err)
				// Если произошла ошибка соединения, допустим упал клиент, то убираем это соединение.
				client.doneWrite <- true
				return
			} else {
				AppServer.IncomingMessage(client, &msg)
			}
		}
	}
}
