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

func NewClient(ws *websocket.Conn) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	ch := make(chan string, channalBufSize)
	doneCh := make(chan bool)
	newClient := &Client{maxID, ws, ch, doneCh}
	maxID++

	return newClient
}
