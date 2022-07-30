package rabbitmq

import (
	"UserApp/src/type/util"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

func Connect(conStr string) *amqp091.Connection {
	conn, err := amqp091.Dial(conStr)
	fmt.Println("conn->", conn)
	fmt.Println("err->", err)
	util.FailOnError(err, "Failed to connect RabbitMq.")
	//defer conn.Close()
	return conn
}

func OpenChannel(conn *amqp091.Connection) *amqp091.Channel {
	channel, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel.")
	//defer channel.Close()
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

func PublishToQueue(channel *amqp091.Channel, queue amqp091.Queue, msg any) bool {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(msg)

	err := channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        reqBodyBytes.Bytes(),
		})

	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()
	//
	//err = channel.PublishWithContext(
	//	ctx,
	//	"",
	//	"",
	//	false,
	//	false,
	//	amqp091.Publishing{
	//		ContentType: "text/plain",
	//		Body:        reqBodyBytes.Bytes(),
	//	})

	util.FailOnError(err, "Failed to publishing data.")
	log.Printf(" [x] Sent %s\n", msg)

	return true
}
