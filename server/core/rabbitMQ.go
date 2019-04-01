package core

import (
	"encoding/json"
	"io"

	"../logger"
	"github.com/streadway/amqp"
)

// var channel amqp.Channel

// Message - Формат сообщений для обмена по RabbitMQ
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
func StartRabbitMQ() {
	//--------------------------------- Owerall ----------------------------
	// Создадим связь с брокером.
	conn, err := amqp.Dial("amqp://macroserv:12345@localhost:5672/macroserv")
	checkError(err, "Failed to connect to RabbitMQ")
	Server.connectRMQ = conn

	// Устанавливаем соединение с брокером.
	channel, err := conn.Channel()
	checkError(err, "Failed to open a channel")
	Server.channelRMQ = channel

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
	checkError(err, "Failed to declare rooms exchange")

	//--------------------------------- For core ----------------------------
	// Создаем очередь из которой будем поулчать сообщения.
	// Делается всегда и там где принимается и там где отправляется,
	// если очереди нет то сообщение просто проигнорится,
	// но если очередь оздана хотя бы раз, то повторно создана не будет.
	// так как это очередь для того что бы слушать сообщения приходящие нам,
	// не надо его запоминать, у нас будет горутина крутится...
	queue, err := channel.QueueDeclare(
		"сore", // name
		true,   // durable
		false,  // delete when usused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	checkError(err, "Failed to declare a queue")

	err = channel.QueueBind(
		queue.Name, // queue name
		queue.Name, // routing key (binding_key)
		// TODO: наверно надо вынести в отдельную переменную.
		"core", // exchange
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

	// Запускаем горутину которая будет "слушать" очередь.
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
				// TODO: Написать тесты для этого метода.
				APIMetods[msg.HandlerName](msg.Data)
			}
		}
	}()
}

// CreateMessage - Запаковывает структуру для отправки.
func CreateMessage(consumerName string, data interface{}, methodName string) {
	message, err := json.Marshal(data)
	if err != nil {
		logger.WarningPrintf("Ошибка при запаковке данных для отправки %v: %v", methodName, err)
		return
	}

	messageRMQ := MessageRMQ{
		HandlerName: methodName,
		Data:        string(message),
	}

	PublishMessage(consumerName, messageRMQ)
}

// PublishMessage - Отправка сообщений в очередь
func PublishMessage(consumerName string, message MessageRMQ) {
	// TODO: Канала может не быть!
	jsonMessag, err := json.Marshal(message)
	checkError(err, "Failed marshal message")

	err = Server.channelRMQ.Publish(
		"rooms",      // exchange
		consumerName, // routing key
		false,        // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(jsonMessag),
		})

	checkError(err, "Failed to publish a message")
}
