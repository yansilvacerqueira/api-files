package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/yansilvacerqueira/api-files/internal/users/entity"
	"github.com/yansilvacerqueira/api-files/internal/users/repository"
)

type Handler struct {
	db     *sql.DB
	logger *log.Logger
	repo   *repository.UserRepository
}

type Config struct {
	DB     *sql.DB
	Logger *log.Logger
}

type createUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func NewHandler(cfg Config) (*Handler, error) {
	if cfg.DB == nil {
		return nil, errors.New("database connection is required")
	}

	logger := cfg.Logger
	if logger == nil {
		logger = log.Default()
	}

	repo := repository.NewUserRepository(cfg.DB)

	return &Handler{
		db:     cfg.DB,
		logger: logger,
		repo:   repo,
	}, nil
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.repo.GetUsers(ctx)
	if err != nil {
		h.logger.Printf("Error fetching users: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	var sanitizedUsers []map[string]interface{}
	for _, user := range users {
		sanitizedUsers = append(sanitizedUsers, user.Sanitize())
	}

	h.respondWithJSON(w, http.StatusOK, sanitizedUsers)
}

func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(path.Base(r.URL.Path), 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	ctx := r.Context()
	user, err := h.repo.GetUserByID(ctx, id)
	if err != nil {
		h.logger.Printf("Error fetching user %d: %v", id, err)
		h.respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	h.respondWithJSON(w, http.StatusOK, user.Sanitize())
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	user, err := entity.NewUser(req.FullName, req.Email, req.Password)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	if err := h.repo.CreateUser(ctx, user); err != nil {
		h.logger.Printf("Error creating user: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, user.Sanitize())
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(path.Base(r.URL.Path), 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	ctx := r.Context()
	user, err := h.repo.GetUserByID(ctx, id)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		if err := entity.ValidateEmail(req.Email); err != nil {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		user.Email = req.Email
	}
	if req.Password != "" {
		if err := user.SetPassword(req.Password); err != nil {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	if err := h.repo.UpdateUser(ctx, user); err != nil {
		h.logger.Printf("Error updating user %d: %v", id, err)
		h.respondWithError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	h.respondWithJSON(w, http.StatusOK, user.Sanitize())
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(path.Base(r.URL.Path), 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	ctx := r.Context()
	if err := h.repo.DeleteUser(ctx, id); err != nil {
		h.logger.Printf("Error deleting user %d: %v", id, err)
		h.respondWithError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
}

func (h *Handler) handleUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (requires authentication middleware)
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx := r.Context()
	user, err := h.repo.GetUserByID(ctx, userID)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	h.respondWithJSON(w, http.StatusOK, user.Sanitize())
}
