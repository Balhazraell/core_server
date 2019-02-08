package websockets

import (
	"encoding/json"
	"io"
	"sync"

	"golang.org/x/net/websocket"

	"../logger"
)

const channalBufSize = 100

var maxID int

// IncomingMessage - Структура описывающая формат входящего сообщения от клиента.
// TODO: Если у отправленных и пришедших сообщений теперь общий вид то надо сделать одну структуру.
type IncomingMessage struct {
	HandlerName string `json:"handler_name"`
	Data        string `json:"data"`
}

// OutcomingMessage - Структура описывающая формат исходящего сообщения для клиента.
type OutcomingMessage struct {
	HandlerName string `json:"handler_name"`
	Data        string `json:"data"`
}

// Client - Структура описывающая способ хранения связи клиента на сервере и соединения с клиентом.
type Client struct {
	id int // Должен браться из базы в соответствии с id пользователя в базе.
	ws *websocket.Conn
	ch chan string // Канал для отправки сообщения.

	shutdownRead chan bool
	// shutdownWrite chan bool
}

// Shutdown - метод завершения соединения с клиентом.
func (client *Client) Shutdown() {
	client.shutdownRead <- true
	// client.shutdownWrite <- true
}

// Listen - метод начала соединения с клиентом. Начинаем слушать сообщения от пользователя.
func (client *Client) Listen() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		client.listenRead()
	}()

	wg.Wait()

	logger.InfoPrintf("Listen у клиента %d работу закончил.", client.id)
	AppServer.DelClient(client)
}

// SetGameMap - отдаем пользователю карту где он сейчас находится.
func (client *Client) SetGameMap(gameMap []byte) {
	// Получили карту - отправляем её пользователю.
	logger.InfoPrintf("Пытаемся задать карту клиенту с id = %v.", client.id)

	newMessage := OutcomingMessage{
		HandlerName: "set_grid",
		Data:        string(gameMap),
	}

	logger.InfoPrint("newMessage сформирован и готовится к отправке.")

	websocket.JSON.Send(client.ws, newMessage)
}

func (client *Client) listenRead() {
	defer func() {
		logger.InfoPrintf("listenRead у клиента %d работу закончил.", client.id)
	}()

	for {
		select {
		case <-client.shutdownRead:
			return
		default:
			var msg IncomingMessage
			err := websocket.JSON.Receive(client.ws, &msg)

			if err == io.EOF {
				return
			} else if err != nil {
				logger.ErrorPrintf("Проблема чтения сообщения от клиента : %v.", err)
				return
			} else {
				AppServer.IncomingMessage(client, &msg)
			}
		}
	}
}

// SendError - отослать пользователю сообщение об ошибке.
func (client *Client) SendError(message string) {
	jsonMessage, err := json.Marshal(message)

	if err != nil {
		logger.ErrorPrintf("При формировнии json при отправки сообщения об ошибке произошла ошибка: %v", err)
		return
	}

	newMessage := OutcomingMessage{
		HandlerName: "send_error",
		Data:        string(jsonMessage),
	}

	websocket.JSON.Send(client.ws, newMessage)
}

// SetRoomsCatalog - метод отдающий пользователю текущий каталок с комнатами.
func (client *Client) SetRoomsCatalog(roomsIDs []int) {
	jsonRoomsIDs, err := json.Marshal(roomsIDs)

	if err != nil {
		logger.ErrorPrintf("При формировнии json при отправки сообщения об ошибке произошла ошибка: %v", err)
		return
	}

	newMessage := OutcomingMessage{
		HandlerName: "set_rooms_catalog",
		Data:        string(jsonRoomsIDs),
	}

	websocket.JSON.Send(client.ws, newMessage)
}
