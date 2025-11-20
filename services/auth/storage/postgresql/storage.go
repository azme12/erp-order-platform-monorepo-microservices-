package postgresql

import (
	"context"
	"database/sql"
	"microservice-challenge/package/errors"
	"microservice-challenge/services/auth/model"
	"microservice-challenge/services/auth/storage/postgresql/db"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Storage struct {
	queries *db.Queries
	db      *sql.DB
}

func NewStorage(database *sql.DB) *Storage {
	return &Storage{
		queries: db.New(database),
		db:      database,
	}
}

// convertDBUserToModel converts sqlc generated db.User to model.User
func convertDBUserToModel(dbUser db.User) *model.User {
	user := &model.User{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Role:         dbUser.Role,
	}

	if dbUser.ResetToken.Valid {
		token := dbUser.ResetToken.String
		user.ResetToken = &token
	}
	if dbUser.ResetTokenExpiresAt.Valid {
		user.ResetTokenExpiresAt = &dbUser.ResetTokenExpiresAt.Time
	}
	if dbUser.CreatedAt.Valid {
		user.CreatedAt = dbUser.CreatedAt.Time
	} else {
		user.CreatedAt = time.Now()
	}
	if dbUser.UpdatedAt.Valid {
		user.UpdatedAt = dbUser.UpdatedAt.Time
	} else {
		user.UpdatedAt = time.Now()
	}

	return user
}

// convertModelUserToCreateParams converts model.User to sqlc CreateUserParams
func convertModelUserToCreateParams(user *model.User) db.CreateUserParams {
	params := db.CreateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	}

	if user.ResetToken != nil {
		params.ResetToken = sql.NullString{
			String: *user.ResetToken,
			Valid:  true,
		}
	}
	if user.ResetTokenExpiresAt != nil {
		params.ResetTokenExpiresAt = sql.NullTime{
			Time:  *user.ResetTokenExpiresAt,
			Valid: true,
		}
	}
	params.CreatedAt = sql.NullTime{
		Time:  user.CreatedAt,
		Valid: !user.CreatedAt.IsZero(),
	}
	params.UpdatedAt = sql.NullTime{
		Time:  user.UpdatedAt,
		Valid: !user.UpdatedAt.IsZero(),
	}

	return params
}

// convertModelUserToUpdateParams converts model.User to sqlc UpdateUserParams
func convertModelUserToUpdateParams(user *model.User) db.UpdateUserParams {
	params := db.UpdateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	}

	if user.ResetToken != nil {
		params.ResetToken = sql.NullString{
			String: *user.ResetToken,
			Valid:  true,
		}
	}
	if user.ResetTokenExpiresAt != nil {
		params.ResetTokenExpiresAt = sql.NullTime{
			Time:  *user.ResetTokenExpiresAt,
			Valid: true,
		}
	}
	params.UpdatedAt = sql.NullTime{
		Time:  user.UpdatedAt,
		Valid: !user.UpdatedAt.IsZero(),
	}

	return params
}

func (s *Storage) CreateUser(ctx context.Context, user *model.User) error {
	params := convertModelUserToCreateParams(user)
	err := s.queries.CreateUser(ctx, params)

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
	dbUser, err := s.queries.GetUserByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	return convertDBUserToModel(dbUser), nil
}

func (s *Storage) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.ErrBadRequest
	}

	dbUser, err := s.queries.GetUserByID(ctx, userID)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	return convertDBUserToModel(dbUser), nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *model.User) error {
	params := convertModelUserToUpdateParams(user)
	err := s.queries.UpdateUser(ctx, params)

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
	dbUser, err := s.queries.GetUserByResetToken(ctx, resetToken)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	return convertDBUserToModel(dbUser), nil
}
