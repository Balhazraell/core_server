package room

import (
	"encoding/json"
	"io"

	"../../logger"
	"github.com/streadway/amqp"
)

var channel amqp.Channel

// MessageRMQ - Формат сообщений для обмена по RabbitMQ
type MessageRMQ struct {
	HandlerName string `json:"handler_name"`
	Data        string `json:"data"`
}

func checkError(err error, message string) {
	if err != nil {
		logger.ErrorPrintf("%s: %s", message, err)
	}
}

// StartRabbitMQ - Запускает создание очередей в RabbitMQ
func StartRabbitMQ(name string) {
	//---------------------------- Overall ----------------------------
	// Сейчас создадим полноценное соединение для RabbitMQ
	conn, err := amqp.Dial("amqp://macroserv:12345@localhost:5672/macroserv")
	checkError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Создаем канал.
	channel, err := conn.Channel()
	checkError(err, "Failed to open a channel")
	defer channel.Close()

	// Точка доступа должна быть создана, до того как создана очередь.
	// так как слать сообщения в несучествующую точку доступа запрещено!
	err = channel.ExchangeDeclare(
		"core",   // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	checkError(err, "Failed to declare core exchange")

	err = channel.ExchangeDeclare(
		"rooms",  // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	checkError(err, "Failed to declare room exchange")

	//---------------------------- For Room ----------------------------
	// Сначала создаем очередь на получение сообщений, назвние
	// будет формироваться из имени комнаты, в нашем случае из id
	queue, err := channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	checkError(err, "Failed to declare a queue")

	// Связываем комнату и точку доступа.
	err = channel.QueueBind(
		queue.Name, // queue name
		queue.Name, // routing key (binding_key)
		"rooms",    // exchange
		false,
		nil,
	)
	checkError(err, "Failed to bind a queue")

	// Теперь создаем подписчика.
	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	checkError(err, "Failed to register a consumer")

	// Мониторим очередь на наличие сообщений.
	go func() {
		for d := range msgs {
			var msg MessageRMQ
			err := json.Unmarshal(d.Body, &msg)

			if err == io.EOF {
				continue
			} else if err != nil {
				logger.ErrorPrintf("Проблема чтения сообщения от комнаты : %v.", err)
				continue
			} else {
				// TODO: можем упасть при вызове не верного метода - надо обработать!
				// Допустим метод которого нет в списке.
				method, ok := APIMetods[msg.HandlerName]
				if ok {
					method(msg.Data)
				} else {
					logger.WarningPrintf("Попытка вызвать API которго нет или к которому нет доступа: %v.", msg.HandlerName)
					continue
				}
			}
		}
	}()
}

// PublishMessage - Отправка сообщений в очередь
func PublishMessage(message MessageRMQ) {
	jsonMessag, err := json.Marshal(message)
	checkError(err, "Failed marshal message")

	err = channel.Publish(
		"core", // exchange
		"сore", // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(jsonMessag),
		})

	checkError(err, "Failed to publish a message")
}
