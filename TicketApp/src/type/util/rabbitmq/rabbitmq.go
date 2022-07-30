package rabbitmq

import (
	"TicketApp/src/type/util"
	"github.com/rabbitmq/amqp091-go"
)

func Connect(conStr string) *amqp091.Connection {
	conn, err := amqp091.Dial(conStr)
	util.FailOnError(err, "Failed to connect RabbitMq.")
	return conn
}

func OpenChannel(conn *amqp091.Connection) *amqp091.Channel {
	channel, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel.")
	return channel
}

func DeclareAQueue(channel *amqp091.Channel, qName string) amqp091.Queue {
	q, err := channel.QueueDeclare(
		qName,
		false,
		false,
		false,
		false,
		nil)
	util.FailOnError(err, "Failed to declare a queue.")
	return q
}

func ConsumeMessage(channel *amqp091.Channel, queue amqp091.Queue) <-chan amqp091.Delivery {
	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	util.FailOnError(err, "Failed to consume.")

	return msgs
}
