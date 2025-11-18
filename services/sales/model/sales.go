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

// String returns the string representation of OrderStatus
func (s OrderStatus) String() string {
	return string(s)
}

// IsValid checks if the OrderStatus is valid
func (s OrderStatus) IsValid() bool {
	return s == OrderStatusDraft || s == OrderStatusConfirmed || s == OrderStatusPaid
}

// IsDraft checks if the order is in draft status
func (s OrderStatus) IsDraft() bool {
	return s == OrderStatusDraft
}

// IsConfirmed checks if the order is confirmed
func (s OrderStatus) IsConfirmed() bool {
	return s == OrderStatusConfirmed
}

// IsPaid checks if the order is paid
func (s OrderStatus) IsPaid() bool {
	return s == OrderStatusPaid
}

type SalesOrder struct {
	// Identifiers
	ID         uuid.UUID `json:"id" db:"id"`
	CustomerID uuid.UUID `json:"customer_id" db:"customer_id"`

	// Business fields
	Status      OrderStatus `json:"status" db:"status"`
	TotalAmount float64     `json:"total_amount" db:"total_amount"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type OrderItem struct {
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

type SalesOrderWithItems struct {
	SalesOrder
	Items []OrderItem `json:"items"`
}

type CreateOrderRequest struct {
	CustomerID uuid.UUID                `json:"customer_id"`
	Items      []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity int       `json:"quantity"`
}

type UpdateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items"`
}
