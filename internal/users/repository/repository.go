package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yansilvacerqueira/api-files/internal/users/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]entity.User, error) {
	query := `
		SELECT id, full_name, email, password, created_at, updated_at, last_login, deleted
		FROM users
		WHERE deleted = false
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.LastLogin,
			&user.Deleted,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `
		SELECT id, full_name, email, password, created_at, updated_at, last_login, deleted
		FROM users
		WHERE id = $1 AND deleted = false
	`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (full_name, email, password, created_at, updated_at, deleted)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		user.FullName,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
		user.Deleted,
	).Scan(&user.ID)

	return err
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET full_name = $1, email = $2, password = $3, updated_at = $4
		WHERE id = $5 AND deleted = false
	`

	result, err := r.db.ExecContext(ctx, query,
		user.FullName,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	query := `
		UPDATE users
		SET deleted = true, updated_at = NOW()
		WHERE id = $1 AND deleted = false
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	query := `
		UPDATE users
		SET last_login = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted = false
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}
