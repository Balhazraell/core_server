package websockets

import (
	"fmt"
	"io"

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
}

func (client *Client) Done() {
	// Нужно убить данный объект и его связи...

	// так как запущены бесконечные циклы для прослушивания каналов,
	// теперь их надо останоить.
	client.doneCh <- true
}

func (client *Client) Listen() {
	go client.listenWrite()
	// Дичайшее уебанство... данный метод не в горутине, что бы клиент не прекратил работу и
	// не завершил себя раньше времени...
	client.listenRead()
}

func (client *Client) listenWrite() {
	for {
		select {
		case <-client.doneCh:
			AppServer.DelClient(client)
			// если данный "слушатель" первым поймал сообщение об удалении,
			// то удаляем и сообщам дальше что бы следующий "слушатель"
			// поймал сообщение и завершил работу
			client.doneCh <- true
			return
		}
	}
}

func (client *Client) listenRead() {
	for {
		select {
		case <-client.doneCh:
			AppServer.DelClient(client)
			// если данный "слушатель" первым поймал сообщение об удалении,
			// то удаляем и сообщам дальше что бы следующий "слушатель"
			// поймал сообщение и завершил работу
			client.doneCh <- true
		default:
			var msg IncomingMessage
			err := websocket.JSON.Receive(client.ws, &msg)

			if err == io.EOF {
				client.doneCh <- true
			} else if err != nil {
				fmt.Println("Проблема чтения сообщения от клиента.")
				fmt.Println(err)
				// Если произошла ошибка соединения, допустим упал клиент, то убираем это соединение.
				client.doneCh <- true
			} else {
				AppServer.IncomingMessage(client, &msg)
			}
		}
	}
}
