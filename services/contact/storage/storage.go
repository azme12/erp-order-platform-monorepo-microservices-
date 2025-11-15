package storage

import (
	"context"
	"microservice-challenge/services/contact/model"
)

type Storage interface {
	CreateCustomer(ctx context.Context, customer model.Customer) error
	GetCustomerByID(ctx context.Context, id string) (model.Customer, error)
	GetCustomerByEmail(ctx context.Context, email string) (model.Customer, error)
	ListCustomers(ctx context.Context, limit, offset int) ([]model.Customer, error)
	UpdateCustomer(ctx context.Context, customer model.Customer) error
	DeleteCustomer(ctx context.Context, id string) error

	CreateVendor(ctx context.Context, vendor model.Vendor) error
	GetVendorByID(ctx context.Context, id string) (model.Vendor, error)
	GetVendorByEmail(ctx context.Context, email string) (model.Vendor, error)
	ListVendors(ctx context.Context, limit, offset int) ([]model.Vendor, error)
	UpdateVendor(ctx context.Context, vendor model.Vendor) error
	DeleteVendor(ctx context.Context, id string) error
}
