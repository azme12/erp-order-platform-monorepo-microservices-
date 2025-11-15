package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/jwt"
	"microservice-challenge/services/auth/model"
	"microservice-challenge/services/auth/storage"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	storage      storage.Storage
	jwtSecret    string
	userExpHours int
}

func NewUsecase(storage storage.Storage, jwtSecret string, userExpHours int) *Usecase {
	return &Usecase{
		storage:      storage,
		jwtSecret:    jwtSecret,
		userExpHours: userExpHours,
	}
}

func (u *Usecase) Register(ctx context.Context, req model.RegisterRequest) (model.User, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	existingUser, err := u.storage.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return model.User{}, errors.ErrConflict
	}
	if err != nil && err != errors.ErrNotFound {
		return model.User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, errors.ErrInternalServerError
	}

	user := model.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.storage.CreateUser(ctx, &user); err != nil {
		return model.User{}, err
	}

	user.PasswordHash = ""
	return user, nil
}

func (u *Usecase) Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := u.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if err == errors.ErrNotFound {
			return model.LoginResponse{}, errors.ErrUnauthorized
		}
		return model.LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return model.LoginResponse{}, errors.ErrUnauthorized
	}

	token, err := jwt.GenerateUserToken(user.ID.String(), user.Email, user.Role, u.jwtSecret, u.userExpHours)
	if err != nil {
		return model.LoginResponse{}, errors.ErrInternalServerError
	}

	user.PasswordHash = ""

	const secondsPerHour = 3600
	return model.LoginResponse{
		AccessToken: token,
		ExpiresIn:   u.userExpHours * secondsPerHour,
		User:        *user,
	}, nil
}

func (u *Usecase) ForgotPassword(ctx context.Context, req model.ForgotPasswordRequest) (model.ForgotPasswordResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := u.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if err == errors.ErrNotFound {
			return model.ForgotPasswordResponse{
				Message: "If the email exists, a password reset link has been sent",
			}, nil
		}
		return model.ForgotPasswordResponse{}, err
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return model.ForgotPasswordResponse{}, errors.ErrInternalServerError
	}
	resetToken := hex.EncodeToString(tokenBytes)

	expiresAt := time.Now().Add(1 * time.Hour)

	user.ResetToken = &resetToken
	user.ResetTokenExpiresAt = &expiresAt
	user.UpdatedAt = time.Now()

	if err := u.storage.UpdateUser(ctx, user); err != nil {
		return model.ForgotPasswordResponse{}, err
	}

	return model.ForgotPasswordResponse{
		Message:    "If the email exists, a password reset link has been sent",
		ResetToken: resetToken,
	}, nil
}

func (u *Usecase) ResetPassword(ctx context.Context, req model.ResetPasswordRequest) error {
	user, err := u.storage.GetUserByResetToken(ctx, req.ResetToken)
	if err != nil {
		if err == errors.ErrNotFound {
			return errors.ErrInvalidToken
		}
		return err
	}

	if user.ResetTokenExpiresAt == nil || time.Now().After(*user.ResetTokenExpiresAt) {
		return errors.ErrTokenExpired
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrInternalServerError
	}

	user.PasswordHash = string(hashedPassword)
	user.ResetToken = nil
	user.ResetTokenExpiresAt = nil
	user.UpdatedAt = time.Now()

	if err := u.storage.UpdateUser(ctx, user); err != nil {
		return err
	}

	return nil
}
