package main

type Server struct {
	// TODO: Непонятно зачем нужен pattern в данном случае.
	pattern   string // Адрес сервера.
	clients map[int]
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

func DelClient(client *Client) {
	delete(Server.clients, client.id)
}