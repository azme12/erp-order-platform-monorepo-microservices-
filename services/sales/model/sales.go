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

type SalesOrder struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	CustomerID  uuid.UUID   `json:"customer_id" db:"customer_id"`
	Status      OrderStatus `json:"status" db:"status"`
	TotalAmount float64     `json:"total_amount" db:"total_amount"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

type OrderItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	OrderID   uuid.UUID `json:"order_id" db:"order_id"`
	ItemID    uuid.UUID `json:"item_id" db:"item_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	UnitPrice float64   `json:"unit_price" db:"unit_price"`
	Subtotal  float64   `json:"subtotal" db:"subtotal"`
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
