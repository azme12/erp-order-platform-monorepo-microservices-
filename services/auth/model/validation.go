package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (r *RegisterRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 100)),
		validation.Field(&r.Role, validation.Required, validation.In("inventory_manager", "finance_manager")),
	)
}

func (r *LoginRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required),
	)
}

func (r *ForgotPasswordRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required, is.Email),
	)
}

func (r *ResetPasswordRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ResetToken, validation.Required),
		validation.Field(&r.NewPassword, validation.Required, validation.Length(6, 100)),
	)
}
