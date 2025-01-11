package queue

import "encoding/json"

type QueueMessage struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	ID       int    `json:"id"`
}

// Marshal converts QueueMessage to JSON bytes
func (m *QueueMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal parses JSON bytes into QueueMessage
func (m *QueueMessage) FromJSON(data []byte) error {
	return json.Unmarshal(data, m)
}
