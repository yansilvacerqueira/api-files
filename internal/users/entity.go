package users

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidPassword = errors.New("password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	ErrEmptyFullName   = errors.New("full name is required")
	ErrEmptyPassword   = errors.New("password is required")
)

type User struct {
	ID        int64
	FullName  string
	Email     string
	Password  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   bool
	LastLogin *time.Time
}

// UserStatus represents the current state of a user
type UserStatus struct {
	IsActive    bool
	IsLocked    bool
	LastLoginIP string
	LoginCount  int
}

func NewUser(fullName, email, password string) (*User, error) {
	now := time.Now()

	// Validate and sanitize inputs
	fullName = strings.TrimSpace(fullName)
	email = strings.TrimSpace(strings.ToLower(email))

	if fullName == "" {
		return nil, ErrEmptyFullName
	}

	if err := validateEmail(email); err != nil {
		return nil, err
	}

	u := &User{
		FullName:  fullName,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.SetPassword(password); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) SetPassword(password string) error {
	if password == "" {
		return ErrEmptyPassword
	}

	if err := validatePassword(password); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))

	return err == nil
}

func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}

func (u *User) SoftDelete() {
	u.Deleted = true
	u.UpdatedAt = time.Now()
}

func (u *User) IsDeleted() bool {
	return u.Deleted
}

func validateEmail(email string) error {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return ErrInvalidEmail
	}
	// TODO: refactor that using something like regex
	return nil
}

// validatePassword ensures password meets security requirements
func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrInvalidPassword
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !(hasUpper && hasLower && hasNumber && hasSpecial) {
		return ErrInvalidPassword
	}

	return nil
}

// Sanitize removes sensitive data for safe transmission
func (u *User) Sanitize() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"full_name":  u.FullName,
		"email":      u.Email,
		"created_at": u.CreatedAt,
		"last_login": u.LastLogin,
	}
}
