package users

import (
	"net/http"
)

func (h *Handler) SetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/users", h.handleUsers)
	mux.HandleFunc("/api/users/", h.handleUserByID)
}

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		h.respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUsers(w, r)
	case http.MethodPut:
		h.updateUser(w, r)
	case http.MethodDelete:
		h.deleteUser(w, r)
	default:
		h.respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
