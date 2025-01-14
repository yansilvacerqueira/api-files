package utils

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON envia uma resposta em formato JSON
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// RespondWithError envia uma mensagem de erro
func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, map[string]string{"error": message})
}

// ParseJSONBody faz o parse do body JSON de uma requisição
func ParseJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&dst)
}
