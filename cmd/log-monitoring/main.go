package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	amqpURI := "amqp://admin:admin@localhost:5672/"
	exchange := "bitcoinrateapp-logs"
	queueName := "logs-queue"
	routingKey := "logs"

	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %s", err)
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare queue: %s", err)
	}

	err = channel.QueueBind(
		queueName,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to bind queue: %s", err)
	}

	messages, err := channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to consume messages: %s", err)
	}

	log.Println("Waiting for messages. Press Ctrl+C to exit.")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg := <-messages:
			log.Printf("[%s] Received message: %s\n", time.Now().Format(time.RFC3339), msg.Body)
		case <-signals:
			log.Println("Exiting...")
			return
		}
	}
}
