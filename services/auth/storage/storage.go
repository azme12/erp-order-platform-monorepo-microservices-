package storage

import (
	"context"
	"microservice-challenge/services/auth/model"
)

type Storage interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	GetUserByResetToken(ctx context.Context, resetToken string) (*model.User, error)
}
