package model

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r *CreatePurchaseOrderRequest) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.VendorID, validation.Required),
		validation.Field(&r.Items, validation.Required, validation.Length(1, 100)),
	); err != nil {
		return err
	}

	// Validate each item in the items slice
	for i, item := range r.Items {
		if err := item.Validate(); err != nil {
			return validation.NewError("items", fmt.Sprintf("item[%d]: %v", i, err))
		}
	}

	return nil
}

func (r *CreatePurchaseOrderItemRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ItemID, validation.Required),
		validation.Field(&r.Quantity, validation.Required, validation.Min(1)),
	)
}

func (r *UpdatePurchaseOrderRequest) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Items, validation.Required, validation.Length(1, 100)),
	); err != nil {
		return err
	}

	// Validate each item in the items slice
	for i, item := range r.Items {
		if err := item.Validate(); err != nil {
			return validation.NewError("items", fmt.Sprintf("item[%d]: %v", i, err))
		}
	}

	return nil
}
