package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (r *CreateCustomerRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Phone, validation.Length(0, 50)),
		validation.Field(&r.Address, validation.Length(0, 500)),
	)
}

func (r *UpdateCustomerRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Phone, validation.Length(0, 50)),
		validation.Field(&r.Address, validation.Length(0, 500)),
	)
}

func (r *CreateVendorRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Phone, validation.Length(0, 50)),
		validation.Field(&r.Address, validation.Length(0, 500)),
	)
}

func (r *UpdateVendorRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Phone, validation.Length(0, 50)),
		validation.Field(&r.Address, validation.Length(0, 500)),
	)
}
