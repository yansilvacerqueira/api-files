package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL       string
	QueueName string
	Timeout   time.Duration
}

type RabbitMQConnection struct {
	config RabbitMQConfig
	conn   *amqp091.Connection
}

// PublishMessage sends a message to the RabbitMQ queue
func (rc *RabbitMQConnection) PublishMessage(msg []byte) error {
	channel, err := rc.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	defer channel.Close()

	messageProperties := amqp091.Publishing{
		DeliveryMode: amqp091.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         msg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return channel.PublishWithContext(ctx, "", rc.config.QueueName, false, false, messageProperties)
}

// ReceiveMessage listens for messages from the RabbitMQ queue
func (rc *RabbitMQConnection) ReceiveMessage(c chan<- QueueMessage) error {
	channel, err := rc.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		rc.config.QueueName,
		false, // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	messages, err := channel.Consume(
		queue.Name,
		"",    // consumer
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	for msg := range messages {
		var queueMessage QueueMessage

		if err := queueMessage.FromJSON(msg.Body); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		c <- queueMessage
	}

	return nil
}

// createRabbitMQConnection initializes a new RabbitMQ connection
func createRabbitMQConnection(cfg RabbitMQConfig) (*RabbitMQConnection, error) {
	conn, err := amqp091.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return &RabbitMQConnection{
		config: cfg,
		conn:   conn,
	}, nil
}
