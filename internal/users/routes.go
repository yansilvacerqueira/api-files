package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// TODO: finish implementation

// Handler manages user-related HTTP handlers with a shared DB connection
type Handler struct {
	db     *sql.DB
	logger *log.Logger
}

// Config holds handler configuration
type Config struct {
	DB     *sql.DB
	Logger *log.Logger
}

// Response represents a standardized HTTP response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewHandler creates a new user handler with dependencies
func NewHandler(cfg Config) (*Handler, error) {
	if cfg.DB == nil {
		return nil, errors.New("database connection is required")
	}

	// Use provided logger or create a default one
	logger := cfg.Logger
	if logger == nil {
		logger = log.Default()
	}

	return &Handler{
		db:     cfg.DB,
		logger: logger,
	}, nil
}

// SetRoutes registers all user-related routes
func (h *Handler) SetRoutes(mux *http.ServeMux) {
	// User management endpoints
	mux.HandleFunc("/api/users", h.handleUsers)
	mux.HandleFunc("/api/users/", h.handleUserByID) // For paths like /api/users/123
	mux.HandleFunc("/api/users/profile", h.handleUserProfile)
}

// Middleware for checking database health
func (h *Handler) withDBCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if database is accessible
		ctx, cancel := context.WithTimeout(r.Context(), 1)
		defer cancel()

		if err := h.db.PingContext(ctx); err != nil {
			h.respondWithError(w, http.StatusServiceUnavailable, "database is not available")
			return
		}
		next(w, r)
	}
}

// handleUsers handles GET and POST requests for /api/users
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

// handleUserByID handles requests for specific users
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

// Utility methods for consistent response handling
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

// Example implementation of a handler method
func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.getUsersFromDB(ctx)
	if err != nil {
		h.logger.Printf("Error fetching users: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	// Sanitize user data before sending
	var sanitizedUsers []map[string]interface{}
	for _, user := range users {
		sanitizedUsers = append(sanitizedUsers, user.Sanitize())
	}

	h.respondWithJSON(w, http.StatusOK, sanitizedUsers)
}

// getUsersFromDB fetches users from the database
func (h *Handler) getUsersFromDB(ctx context.Context) ([]User, error) {
	query := `
		SELECT id, full_name, email, created_at, updated_at, last_login
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := h.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.LastLogin,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}
