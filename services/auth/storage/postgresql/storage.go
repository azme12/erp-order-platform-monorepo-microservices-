package postgresql

import (
	"context"
	"database/sql"
	"microservice-challenge/package/errors"
	"microservice-challenge/services/auth/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.ErrConflict
			}
		}
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, role, reset_token, reset_token_expires_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &model.User{}
	var resetToken sql.NullString
	var resetTokenExpiresAt sql.NullTime
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&resetToken,
		&resetTokenExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	if resetToken.Valid {
		token := resetToken.String
		user.ResetToken = &token
	}
	if resetTokenExpiresAt.Valid {
		user.ResetTokenExpiresAt = &resetTokenExpiresAt.Time
	}

	return user, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.ErrBadRequest
	}

	query := `
		SELECT id, email, password_hash, role, reset_token, reset_token_expires_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &model.User{}
	var resetToken sql.NullString
	var resetTokenExpiresAt sql.NullTime
	err = s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&resetToken,
		&resetTokenExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	if resetToken.Valid {
		token := resetToken.String
		user.ResetToken = &token
	}
	if resetTokenExpiresAt.Valid {
		user.ResetTokenExpiresAt = &resetTokenExpiresAt.Time
	}

	return user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET email = $2, password_hash = $3, role = $4, reset_token = $5, reset_token_expires_at = $6, updated_at = $7
		WHERE id = $1
	`

	var resetToken interface{}
	var resetTokenExpiresAt interface{}
	if user.ResetToken != nil {
		resetToken = *user.ResetToken
	} else {
		resetToken = nil
	}
	if user.ResetTokenExpiresAt != nil {
		resetTokenExpiresAt = *user.ResetTokenExpiresAt
	} else {
		resetTokenExpiresAt = nil
	}

	_, err := s.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
		resetToken,
		resetTokenExpiresAt,
		user.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.ErrConflict
			}
		}
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) GetUserByResetToken(ctx context.Context, resetToken string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, role, reset_token, reset_token_expires_at, created_at, updated_at
		FROM users
		WHERE reset_token = $1
	`

	user := &model.User{}
	var token sql.NullString
	var resetTokenExpiresAt sql.NullTime
	err := s.db.QueryRowContext(ctx, query, resetToken).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&token,
		&resetTokenExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	if token.Valid {
		tokenStr := token.String
		user.ResetToken = &tokenStr
	}
	if resetTokenExpiresAt.Valid {
		user.ResetTokenExpiresAt = &resetTokenExpiresAt.Time
	}

	return user, nil
}
