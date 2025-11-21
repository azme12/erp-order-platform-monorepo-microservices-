package model

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID  uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	SKU string    `json:"sku" db:"sku" example:"SKU-001"`

	Name        string  `json:"name" db:"name" example:"Laptop Computer"`
	Description string  `json:"description" db:"description" example:"High-performance laptop with 16GB RAM and 512GB SSD"`
	UnitPrice   float64 `json:"unit_price" db:"unit_price" example:"1299.99"`

	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2025-11-20T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2025-11-20T12:00:00Z"`
}

type Stock struct {
	ID     uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ItemID uuid.UUID `json:"item_id" db:"item_id" example:"550e8400-e29b-41d4-a716-446655440000"`

	Quantity int `json:"quantity" db:"quantity" example:"100"`

	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2025-11-20T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2025-11-20T12:00:00Z"`
}

type CreateItemRequest struct {
	Name        string  `json:"name" example:"Laptop Computer"`
	Description string  `json:"description" example:"High-performance laptop with 16GB RAM and 512GB SSD"`
	SKU         string  `json:"sku" example:"SKU-001"`
	UnitPrice   float64 `json:"unit_price" example:"1299.99"`
}

type UpdateItemRequest struct {
	Name        string  `json:"name" example:"Laptop Computer Updated"`
	Description string  `json:"description" example:"Updated description for laptop"`
	SKU         string  `json:"sku" example:"SKU-001-UPDATED"`
	UnitPrice   float64 `json:"unit_price" example:"1199.99"`
}

type AdjustStockRequest struct {
	Quantity int `json:"quantity" example:"10"`
}
