package model

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`

	Name    string `json:"name" db:"name" example:"John Doe"`
	Email   string `json:"email" db:"email" example:"john.doe@example.com"`
	Phone   string `json:"phone" db:"phone" example:"+251912345678"`
	Address string `json:"address" db:"address" example:"123 Main Street, City, State 12345"`

	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2025-11-20T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2025-11-20T12:00:00Z"`
}

type Vendor struct {
	ID uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`

	Name    string `json:"name" db:"name" example:"Acme Corporation"`
	Email   string `json:"email" db:"email" example:"contact@acme.com"`
	Phone   string `json:"phone" db:"phone" example:"+251955555555"`
	Address string `json:"address" db:"address" example:"999 Business Boulevard, City, State 99999"`

	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2025-11-20T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2025-11-20T12:00:00Z"`
}

type CreateCustomerRequest struct {
	Name    string `json:"name" example:"John Doe"`
	Email   string `json:"email" example:"john.doe@example.com"`
	Phone   string `json:"phone" example:"+251912345678"`
	Address string `json:"address" example:"123 Main Street, City, State 12345"`
}

type UpdateCustomerRequest struct {
	Name    string `json:"name" example:"John Doe Updated"`
	Email   string `json:"email" example:"john.updated@example.com"`
	Phone   string `json:"phone" example:"+251911111111"`
	Address string `json:"address" example:"789 Updated Street, City, State 99999"`
}

type CreateVendorRequest struct {
	Name    string `json:"name" example:"Acme Corporation"`
	Email   string `json:"email" example:"contact@acme.com"`
	Phone   string `json:"phone" example:"+251955555555"`
	Address string `json:"address" example:"999 Business Boulevard, City, State 99999"`
}

type UpdateVendorRequest struct {
	Name    string `json:"name" example:"Acme Corporation Updated"`
	Email   string `json:"email" example:"updated@acme.com"`
	Phone   string `json:"phone" example:"+251966666666"`
	Address string `json:"address" example:"888 Updated Boulevard, City, State 88888"`
}
