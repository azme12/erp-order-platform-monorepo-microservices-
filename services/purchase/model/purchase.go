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
	ID       uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	VendorID uuid.UUID `json:"vendor_id" db:"vendor_id" example:"550e8400-e29b-41d4-a716-446655440001"`

	Status      PurchaseOrderStatus `json:"status" db:"status" example:"Draft"`
	TotalAmount float64             `json:"total_amount" db:"total_amount" example:"2599.98"`

	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2025-11-20T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2025-11-20T12:00:00Z"`
}

type PurchaseOrderItem struct {
	ID      uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440002"`
	OrderID uuid.UUID `json:"order_id" db:"order_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ItemID  uuid.UUID `json:"item_id" db:"item_id" example:"550e8400-e29b-41d4-a716-446655440003"`

	Quantity  int     `json:"quantity" db:"quantity" example:"2"`
	UnitPrice float64 `json:"unit_price" db:"unit_price" example:"1299.99"`
	Subtotal  float64 `json:"subtotal" db:"subtotal" example:"2599.98"`

	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2025-11-20T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2025-11-20T12:00:00Z"`
}

type PurchaseOrderWithItems struct {
	PurchaseOrder
	Items []PurchaseOrderItem `json:"items"`
}

type CreatePurchaseOrderRequest struct {
	VendorID uuid.UUID                        `json:"vendor_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Items    []CreatePurchaseOrderItemRequest `json:"items"`
}

type CreatePurchaseOrderItemRequest struct {
	ItemID   uuid.UUID `json:"item_id" example:"550e8400-e29b-41d4-a716-446655440003"`
	Quantity int       `json:"quantity" example:"2"`
}

type UpdatePurchaseOrderRequest struct {
	Items []CreatePurchaseOrderItemRequest `json:"items"`
}
