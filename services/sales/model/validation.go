package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r *CreateOrderRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.CustomerID, validation.Required),
		validation.Field(&r.Items, validation.Required, validation.Length(1, 100)),
	)
}

func (r *CreateOrderItemRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ItemID, validation.Required),
		validation.Field(&r.Quantity, validation.Required, validation.Min(1)),
	)
}

func (r *UpdateOrderRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Items, validation.Required, validation.Length(1, 100)),
	)
}
