package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "Draft"
	OrderStatusConfirmed OrderStatus = "Confirmed"
	OrderStatusPaid      OrderStatus = "Paid"
)

func (s OrderStatus) String() string {
	return string(s)
}

func (s OrderStatus) IsValid() bool {
	return s == OrderStatusDraft || s == OrderStatusConfirmed || s == OrderStatusPaid
}

func (s OrderStatus) IsDraft() bool {
	return s == OrderStatusDraft
}

func (s OrderStatus) IsConfirmed() bool {
	return s == OrderStatusConfirmed
}

func (s OrderStatus) IsPaid() bool {
	return s == OrderStatusPaid
}

type SalesOrder struct {
	ID         uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CustomerID uuid.UUID `json:"customer_id" db:"customer_id" example:"550e8400-e29b-41d4-a716-446655440001"`

	Status      OrderStatus `json:"status" db:"status" example:"Draft"`
	TotalAmount float64     `json:"total_amount" db:"total_amount" example:"2599.98"`

	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2025-11-20T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2025-11-20T12:00:00Z"`
}

type OrderItem struct {
	ID      uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440002"`
	OrderID uuid.UUID `json:"order_id" db:"order_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ItemID  uuid.UUID `json:"item_id" db:"item_id" example:"550e8400-e29b-41d4-a716-446655440003"`

	Quantity  int     `json:"quantity" db:"quantity" example:"2"`
	UnitPrice float64 `json:"unit_price" db:"unit_price" example:"1299.99"`
	Subtotal  float64 `json:"subtotal" db:"subtotal" example:"2599.98"`

	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2025-11-20T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2025-11-20T12:00:00Z"`
}

type SalesOrderWithItems struct {
	SalesOrder
	Items []OrderItem `json:"items"`
}

type CreateOrderRequest struct {
	CustomerID uuid.UUID                `json:"customer_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Items      []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	ItemID   uuid.UUID `json:"item_id" example:"550e8400-e29b-41d4-a716-446655440003"`
	Quantity int       `json:"quantity" example:"2"`
}

type UpdateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items"`
}
