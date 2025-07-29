package messaging

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"orderworker/utils"
	"os"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
)

func Connect() {
	rabbitMqUrl := os.Getenv("RABBIT_URL")
	if rabbitMqUrl == "" {
		log.Fatal("The RABBITMQ_URL environment variable must be defined.")
	}
	var err error

	conn, err := amqp.Dial(rabbitMqUrl)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err = conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")

	_, err = ch.QueueDeclare(
		"order_events",
		true, false, false, false, nil,
	)
	utils.FailOnError(err, "Failed to declare a queue")

	log.Println("Successfully Connected to RabbitMQ")
}

func Consume(handler func(body []byte)) {
	msgs, err := ch.Consume(
		"order_events",
		"",
		true,
		false,
		false,
		false,
		nil)
	utils.FailOnError(err, "Failed to register a consumer")

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func Close() {
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}
	log.Println("Successfully Disconnected from RabbitMQ")
}
