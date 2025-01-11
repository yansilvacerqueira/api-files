package queue

import "encoding/json"

// QueueDto represents the message structure
type QueueDto struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	ID       int    `json:"id"`
}

// Marshal converts QueueDto to JSON bytes
func (q *QueueDto) Marshal() ([]byte, error) {
	return json.Marshal(q)
}

// Unmarshal parses JSON bytes into QueueDto
func (q *QueueDto) Unmarshal(v []byte) error {
	return json.Unmarshal(v, q)
}
