package main

import (
	"golang.org/x/net/websocket"
)

const channalBufSize = 100

var maxID int

type RoutingMessage struct {
	HandlerName string `json:"handler_name"`
	Data struct `json:"data"`
}

type Client struct {
	id int
	ws *websocket.Conn
	ch chan string
	doneCh chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	ch := make(chan string, channelBufSize)
	doneCh := make(chan bool)
}