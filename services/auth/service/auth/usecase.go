package service

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

const (
	secondsPerHour     = 3600
	resetTokenSize     = 32
	resetTokenValidity = 1 * time.Hour
)

type Service struct {
	storage      storage.Storage
	jwtSecret    string
	userExpHours int
}

func NewService(storage storage.Storage, jwtSecret string, userExpHours int) *Service {
	return &Service{
		storage:      storage,
		jwtSecret:    jwtSecret,
		userExpHours: userExpHours,
	}
}

func (s *Service) Register(ctx context.Context, req model.RegisterRequest) (model.User, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	_, err := s.storage.GetUserByEmail(ctx, email)
	if err == nil {
		return model.User{}, errors.ErrConflict
	}
	if err != errors.ErrNotFound {
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

	if err := s.storage.CreateUser(ctx, &user); err != nil {
		return model.User{}, err
	}

	user.PasswordHash = ""
	return user, nil
}

func (s *Service) Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.storage.GetUserByEmail(ctx, email)
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

	token, err := jwt.GenerateUserToken(user.ID.String(), user.Email, user.Role, s.jwtSecret, s.userExpHours)
	if err != nil {
		return model.LoginResponse{}, errors.ErrInternalServerError
	}

	user.PasswordHash = ""

	return model.LoginResponse{
		AccessToken: token,
		ExpiresIn:   s.userExpHours * secondsPerHour,
		User:        *user,
	}, nil
}

func (s *Service) ForgotPassword(ctx context.Context, req model.ForgotPasswordRequest) (model.ForgotPasswordResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if err == errors.ErrNotFound {
			return model.ForgotPasswordResponse{
				Message: "If the email exists, a password reset link has been sent",
			}, nil
		}
		return model.ForgotPasswordResponse{}, err
	}

	tokenBytes := make([]byte, resetTokenSize)
	if _, err := rand.Read(tokenBytes); err != nil {
		return model.ForgotPasswordResponse{}, errors.ErrInternalServerError
	}
	resetToken := hex.EncodeToString(tokenBytes)

	expiresAt := time.Now().Add(resetTokenValidity)

	user.ResetToken = &resetToken
	user.ResetTokenExpiresAt = &expiresAt
	user.UpdatedAt = time.Now()

	if err := s.storage.UpdateUser(ctx, user); err != nil {
		return model.ForgotPasswordResponse{}, err
	}

	return model.ForgotPasswordResponse{
		Message:    "If the email exists, a password reset link has been sent",
		ResetToken: resetToken,
	}, nil
}

func (s *Service) ResetPassword(ctx context.Context, req model.ResetPasswordRequest) error {
	user, err := s.storage.GetUserByResetToken(ctx, req.ResetToken)
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

	if err := s.storage.UpdateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) GenerateServiceToken(ctx context.Context, req model.ServiceTokenRequest) (model.ServiceTokenResponse, error) {
	expectedSecrets := map[string]string{
		"sales":     s.jwtSecret + "_sales",
		"purchase":  s.jwtSecret + "_purchase",
		"contact":   s.jwtSecret + "_contact",
		"inventory": s.jwtSecret + "_inventory",
	}

	expectedSecret, ok := expectedSecrets[req.ServiceName]
	if !ok || req.ServiceSecret != expectedSecret {
		return model.ServiceTokenResponse{}, errors.ErrUnauthorized
	}

	serviceExpHours := 1
	token, err := jwt.GenerateServiceToken(req.ServiceName, s.jwtSecret, serviceExpHours)
	if err != nil {
		return model.ServiceTokenResponse{}, errors.ErrInternalServerError
	}

	return model.ServiceTokenResponse{
		Token:     token,
		ExpiresIn: serviceExpHours * secondsPerHour,
	}, nil
}
