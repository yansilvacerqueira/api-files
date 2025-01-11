package queue

import (
	"fmt"
	"reflect"
)

// QueueType defines supported queue implementations
type QueueImplementation int

const (
	RabbitMQ QueueImplementation = iota
)

// QueueConnection defines the interface for queue operations
type QueueOperations interface {
	PublishMessage([]byte) error
	ReceiveMessage(chan<- QueueMessage) error
}

// Queue encapsulates a specific queue connection implementation
type Queue struct {
	connection QueueOperations
}

// NewQueue creates a new Queue instance based on the specified implementation type
func NewQueue(queueType QueueImplementation, config any) (*Queue, error) {
	var queue Queue

	configType := reflect.TypeOf(config)

	switch queueType {
	case RabbitMQ:
		if configType.Name() != "RabbitMQConfig" {
			return nil, fmt.Errorf("config must be of type RabbitMQConfig")
		}

		connection, err := createRabbitMQConnection(config.(RabbitMQConfig))
		if err != nil {
			return nil, fmt.Errorf("failed to create RabbitMQ connection: %w", err)
		}

		queue.connection = connection

	default:
		return nil, fmt.Errorf("unsupported queue type")
	}

	return &queue, nil
}

// PublishMessage sends a message to the queue
func (q *Queue) PublishMessage(msg []byte) error {
	if q.connection == nil {
		return fmt.Errorf("queue connection not initialized")
	}

	return q.connection.PublishMessage(msg)
}

// ReceiveMessage reads messages from the queue and sends them to the provided channel
func (q *Queue) ReceiveMessage(c chan<- QueueMessage) error {
	if q.connection == nil {
		return fmt.Errorf("queue connection not initialized")
	}

	return q.connection.ReceiveMessage(c)
}
