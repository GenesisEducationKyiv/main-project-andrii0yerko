package logger

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQLogger struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	exchange   string
	routingKey string
	queueName  string
}

func NewRabbitMQLogger(amqpURI, exchange, routingKey, queueName string) (*RabbitMQLogger, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	logger := &RabbitMQLogger{
		conn:       conn,
		channel:    channel,
		exchange:   exchange,
		routingKey: routingKey,
		queueName:  queueName,
	}

	// Declare the channel and queue on-the-fly
	err = logger.declareChannelAndQueue()
	if err != nil {
		logger.Close()
		return nil, err
	}

	return logger, nil
}

func (l *RabbitMQLogger) logMessage(level string, message string) {
	logMessage := fmt.Sprintf("[%s] %s - %s", level, time.Now().Format(time.RFC3339), message)
	err := l.channel.Publish(
		l.exchange,
		l.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(logMessage),
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %s\n", err.Error())
	}
}

func (l *RabbitMQLogger) Debug(message string) {
	l.logMessage(logLevelDebug, message)
}

func (l *RabbitMQLogger) Info(message string) {
	l.logMessage(logLevelInfo, message)
}

func (l *RabbitMQLogger) Error(message string) {
	l.logMessage(logLevelError, message)
}

func (l *RabbitMQLogger) Close() {
	if l.channel != nil {
		if err := l.channel.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %s\n", err.Error())
		}
	}

	if l.conn != nil {
		if err := l.conn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %s\n", err.Error())
		}
	}
}

func (l *RabbitMQLogger) declareChannelAndQueue() error {
	// Declare the exchange
	err := l.channel.ExchangeDeclare(
		l.exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare the queue
	_, err = l.channel.QueueDeclare(
		l.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Bind the queue to the exchange with the routing key
	err = l.channel.QueueBind(
		l.queueName,
		l.routingKey,
		l.exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
