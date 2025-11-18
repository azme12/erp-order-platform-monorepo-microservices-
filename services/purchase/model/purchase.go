package model

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrderStatus string

const (
	PurchaseOrderStatusDraft    PurchaseOrderStatus = "Draft"
	PurchaseOrderStatusReceived PurchaseOrderStatus = "Received"
	PurchaseOrderStatusPaid     PurchaseOrderStatus = "Paid"
)

type PurchaseOrder struct {
	// Identifiers
	ID       uuid.UUID `json:"id" db:"id"`
	VendorID uuid.UUID `json:"vendor_id" db:"vendor_id"`

	// Business fields
	Status      PurchaseOrderStatus `json:"status" db:"status"`
	TotalAmount float64             `json:"total_amount" db:"total_amount"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PurchaseOrderItem struct {
	// Identifiers
	ID      uuid.UUID `json:"id" db:"id"`
	OrderID uuid.UUID `json:"order_id" db:"order_id"`
	ItemID  uuid.UUID `json:"item_id" db:"item_id"`

	// Business fields
	Quantity  int     `json:"quantity" db:"quantity"`
	UnitPrice float64 `json:"unit_price" db:"unit_price"`
	Subtotal  float64 `json:"subtotal" db:"subtotal"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PurchaseOrderWithItems struct {
	PurchaseOrder
	Items []PurchaseOrderItem `json:"items"`
}

type CreatePurchaseOrderRequest struct {
	VendorID uuid.UUID                        `json:"vendor_id"`
	Items    []CreatePurchaseOrderItemRequest `json:"items"`
}

type CreatePurchaseOrderItemRequest struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity int       `json:"quantity"`
}

type UpdatePurchaseOrderRequest struct {
	Items []CreatePurchaseOrderItemRequest `json:"items"`
}
