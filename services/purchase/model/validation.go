package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r *CreatePurchaseOrderRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.VendorID, validation.Required),
		validation.Field(&r.Items, validation.Required, validation.Length(1, 100)),
	)
}

func (r *CreatePurchaseOrderItemRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ItemID, validation.Required),
		validation.Field(&r.Quantity, validation.Required, validation.Min(1)),
	)
}

func (r *UpdatePurchaseOrderRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Items, validation.Required, validation.Length(1, 100)),
	)
}
