package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email string    `json:"email" db:"email" example:"user@example.com"`

	PasswordHash        string     `json:"-" db:"password_hash"`
	Role                string     `json:"role" db:"role" example:"finance_manager"`
	ResetToken          *string    `json:"-" db:"reset_token"`
	ResetTokenExpiresAt *time.Time `json:"-" db:"reset_token_expires_at"`

	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2025-11-20T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2025-11-20T12:00:00Z"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
	Role     string `json:"role" validate:"required,oneof=inventory_manager finance_manager" example:"finance_manager"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNTUwZTg0MDAtZTI5Yi00MWQ0LWE3MTYtNDQ2NjU1NDQwMDAwIiwiZW1haWwiOiJ1c2VyQGV4YW1wbGUuY29tIiwicm9sZSI6ImZpbmFuY2VfbWFuYWdlciIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxNzYzNzI5NjUyLCJpYXQiOjE3NjM2NDMyNTJ9.dxKJuAaIxwQXZa0of9wUREEHsZupFoZ4kWfjnlRrP10"`
	ExpiresIn   int    `json:"expires_in" example:"86400"`
	User        User   `json:"user"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
}

type ForgotPasswordResponse struct {
	Message    string `json:"message" example:"If the email exists, a password reset link has been sent"`
	ResetToken string `json:"reset_token" example:"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6"`
}

type ResetPasswordRequest struct {
	ResetToken  string `json:"reset_token" validate:"required" example:"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6"`
	NewPassword string `json:"new_password" validate:"required,min=6" example:"newpassword123"`
}

type ServiceTokenRequest struct {
	ServiceName   string `json:"service_name" validate:"required" example:"sales"`
	ServiceSecret string `json:"service_secret" validate:"required" example:"your-secret-key-change-in-production_sales"`
}

type ServiceTokenResponse struct {
	Token     string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoic2FsZXMiLCJlbWFpbCI6InNhbGVzQHNlcnZpY2UiLCJyb2xlIjoic2VydmljZSIsInR5cGUiOiJzZXJ2aWNlIiwiZXhwIjoxNzYzNjQ2ODUzLCJpYXQiOjE3NjM2NDMyNTN9.hx7y9wLIVgGoJOhHiSIy7TY91zz6is-NeWeggF-wh1I"`
	ExpiresIn int    `json:"expires_in" example:"3600"`
}
