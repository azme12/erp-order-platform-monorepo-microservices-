package model

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	SKU         string    `json:"sku" db:"sku"`
	UnitPrice   float64   `json:"unit_price" db:"unit_price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Stock struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ItemID    uuid.UUID `json:"item_id" db:"item_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	UnitPrice   float64 `json:"unit_price"`
}

type UpdateItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	UnitPrice   float64 `json:"unit_price"`
}

type AdjustStockRequest struct {
	Quantity int `json:"quantity"`
}
