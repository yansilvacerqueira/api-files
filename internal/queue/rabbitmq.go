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
	TopicName string
	Timeout   time.Time
}

type RabbitConnection struct {
	cfg  RabbitMQConfig
	conn *amqp091.Connection
}

// Publish implements message publishing for RabbitMQ
func (rc *RabbitConnection) Publish(msg []byte) error {
	c, err := rc.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	defer c.Close()

	mp := amqp091.Publishing{
		DeliveryMode: amqp091.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         msg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return c.PublishWithContext(ctx, "", rc.cfg.TopicName, false, false, mp)
}

// Consume implements message consumption for RabbitMQ
func (rc *RabbitConnection) Consume(c chan<- QueueDto) error {
	ch, err := rc.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		rc.cfg.TopicName,
		false, // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name,
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

	for d := range msgs {
		var dto QueueDto

		if err := dto.Unmarshal(d.Body); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		c <- dto
	}

	return nil
}

// newRabbitConn creates a new RabbitMQ connection
func newRabbitConn(cfg RabbitMQConfig) (*RabbitConnection, error) {
	conn, err := amqp091.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return &RabbitConnection{
		cfg:  cfg,
		conn: conn,
	}, nil
}
