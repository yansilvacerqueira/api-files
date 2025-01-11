package queue

import (
	"fmt"
	"reflect"
)

// QueueType defines supported queue implementations
type QueueType int

const (
	RabbitMQ QueueType = iota // Currently only RabbitMQ is supported
)

// QueueConnection defines the interface for queue operations
type QueueConnection interface {
	Publish([]byte) error
	Consume(chan<- QueueDto) error
}

// Queue wraps a queue connection implementation
type Queue struct {
	qc QueueConnection
}

// New creates a new Queue instance based on the specified type
func New(qt QueueType, cfg any) (*Queue, error) {
	var q Queue

	rt := reflect.TypeOf(cfg)

	switch qt {

	case RabbitMQ:
		if rt.Name() != "RabbitMQConfig" {
			return nil, fmt.Errorf("config must be of type RabbitMQConfig")
		}

		conn, err := newRabbitConn(cfg.(RabbitMQConfig))
		if err != nil {
			return nil, fmt.Errorf("failed to create RabbitMQ connection: %w", err)
		}

		q.qc = conn

	default:
		return nil, fmt.Errorf("unsupported queue type")
	}

	return &q, nil
}

// Publish sends a message to the queue
func (q *Queue) Publish(msg []byte) error {
	if q.qc == nil {
		return fmt.Errorf("queue connection not initialized")
	}

	return q.qc.Publish(msg)
}

// Consume reads messages from the queue and sends them to the provided channel
func (q *Queue) Consume(c chan<- QueueDto) error {
	if q.qc == nil {
		return fmt.Errorf("queue connection not initialized")
	}

	return q.qc.Consume(c)
}
