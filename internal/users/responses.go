package users

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response := Response{
		Success: code >= 200 && code < 300,
		Data:    payload,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Error encoding response: %v", err)
	}
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	response := Response{
		Success: false,
		Error:   message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Error encoding error response: %v", err)
	}
}
