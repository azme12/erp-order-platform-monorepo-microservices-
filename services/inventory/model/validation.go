package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r *CreateItemRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Description, validation.Length(0, 1000)),
		validation.Field(&r.SKU, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.UnitPrice, validation.Required, validation.Min(0.0)),
	)
}

func (r *UpdateItemRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Description, validation.Length(0, 1000)),
		validation.Field(&r.SKU, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.UnitPrice, validation.Required, validation.Min(0.0)),
	)
}

func (r *AdjustStockRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Quantity, validation.Required),
	)
}
