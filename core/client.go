package main

import (
	"golang.org/x/net/websocket"
)

const channalBufSize = 100

var maxID int

// type RoutingMessage struct {
// 	HandlerName string `json:"handler_name"`
// 	Data struct `json:"data"`
// }

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

func (client *Client) listenWrite(){
	for{
		select{
		case <- client.doneCh:
			DelClient(client)
			// если данный "слушатель" первым поймал сообщение об удалении, 
			// то удаляем и сообщам дальше что бы следующий "слушатель" 
			// поймал сообщение и завершил работу
			client.doneCh <- true 
			return
		}
	}
}

