package postgresql

import (
	"context"
	"database/sql"
	"microservice-challenge/package/errors"
	"microservice-challenge/services/contact/model"
	"microservice-challenge/services/contact/storage/postgresql/db"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Storage struct {
	queries *db.Queries
	db      *sql.DB
}

func NewStorage(database *sql.DB) *Storage {
	return &Storage{
		queries: db.New(database),
		db:      database,
	}
}

// convertDBCustomerToModel converts sqlc generated db.Customer to model.Customer
func convertDBCustomerToModel(dbCustomer db.Customer) model.Customer {
	customer := model.Customer{
		ID:    dbCustomer.ID,
		Name:  dbCustomer.Name,
		Email: dbCustomer.Email,
	}

	if dbCustomer.Phone.Valid {
		customer.Phone = dbCustomer.Phone.String
	}
	if dbCustomer.Address.Valid {
		customer.Address = dbCustomer.Address.String
	}
	if dbCustomer.CreatedAt.Valid {
		customer.CreatedAt = dbCustomer.CreatedAt.Time
	} else {
		customer.CreatedAt = time.Now()
	}
	if dbCustomer.UpdatedAt.Valid {
		customer.UpdatedAt = dbCustomer.UpdatedAt.Time
	} else {
		customer.UpdatedAt = time.Now()
	}

	return customer
}

// convertModelCustomerToCreateParams converts model.Customer to sqlc CreateCustomerParams
func convertModelCustomerToCreateParams(customer model.Customer) db.CreateCustomerParams {
	params := db.CreateCustomerParams{
		ID:    customer.ID,
		Name:  customer.Name,
		Email: customer.Email,
	}

	if customer.Phone != "" {
		params.Phone = sql.NullString{
			String: customer.Phone,
			Valid:  true,
		}
	}
	if customer.Address != "" {
		params.Address = sql.NullString{
			String: customer.Address,
			Valid:  true,
		}
	}
	params.CreatedAt = sql.NullTime{
		Time:  customer.CreatedAt,
		Valid: !customer.CreatedAt.IsZero(),
	}
	params.UpdatedAt = sql.NullTime{
		Time:  customer.UpdatedAt,
		Valid: !customer.UpdatedAt.IsZero(),
	}

	return params
}

// convertModelCustomerToUpdateParams converts model.Customer to sqlc UpdateCustomerParams
func convertModelCustomerToUpdateParams(customer model.Customer) db.UpdateCustomerParams {
	params := db.UpdateCustomerParams{
		ID:    customer.ID,
		Name:  customer.Name,
		Email: customer.Email,
	}

	if customer.Phone != "" {
		params.Phone = sql.NullString{
			String: customer.Phone,
			Valid:  true,
		}
	}
	if customer.Address != "" {
		params.Address = sql.NullString{
			String: customer.Address,
			Valid:  true,
		}
	}
	params.UpdatedAt = sql.NullTime{
		Time:  customer.UpdatedAt,
		Valid: !customer.UpdatedAt.IsZero(),
	}

	return params
}

// convertDBVendorToModel converts sqlc generated db.Vendor to model.Vendor
func convertDBVendorToModel(dbVendor db.Vendor) model.Vendor {
	vendor := model.Vendor{
		ID:    dbVendor.ID,
		Name:  dbVendor.Name,
		Email: dbVendor.Email,
	}

	if dbVendor.Phone.Valid {
		vendor.Phone = dbVendor.Phone.String
	}
	if dbVendor.Address.Valid {
		vendor.Address = dbVendor.Address.String
	}
	if dbVendor.CreatedAt.Valid {
		vendor.CreatedAt = dbVendor.CreatedAt.Time
	} else {
		vendor.CreatedAt = time.Now()
	}
	if dbVendor.UpdatedAt.Valid {
		vendor.UpdatedAt = dbVendor.UpdatedAt.Time
	} else {
		vendor.UpdatedAt = time.Now()
	}

	return vendor
}

// convertModelVendorToCreateParams converts model.Vendor to sqlc CreateVendorParams
func convertModelVendorToCreateParams(vendor model.Vendor) db.CreateVendorParams {
	params := db.CreateVendorParams{
		ID:    vendor.ID,
		Name:  vendor.Name,
		Email: vendor.Email,
	}

	if vendor.Phone != "" {
		params.Phone = sql.NullString{
			String: vendor.Phone,
			Valid:  true,
		}
	}
	if vendor.Address != "" {
		params.Address = sql.NullString{
			String: vendor.Address,
			Valid:  true,
		}
	}
	params.CreatedAt = sql.NullTime{
		Time:  vendor.CreatedAt,
		Valid: !vendor.CreatedAt.IsZero(),
	}
	params.UpdatedAt = sql.NullTime{
		Time:  vendor.UpdatedAt,
		Valid: !vendor.UpdatedAt.IsZero(),
	}

	return params
}

// convertModelVendorToUpdateParams converts model.Vendor to sqlc UpdateVendorParams
func convertModelVendorToUpdateParams(vendor model.Vendor) db.UpdateVendorParams {
	params := db.UpdateVendorParams{
		ID:    vendor.ID,
		Name:  vendor.Name,
		Email: vendor.Email,
	}

	if vendor.Phone != "" {
		params.Phone = sql.NullString{
			String: vendor.Phone,
			Valid:  true,
		}
	}
	if vendor.Address != "" {
		params.Address = sql.NullString{
			String: vendor.Address,
			Valid:  true,
		}
	}
	params.UpdatedAt = sql.NullTime{
		Time:  vendor.UpdatedAt,
		Valid: !vendor.UpdatedAt.IsZero(),
	}

	return params
}

func (s *Storage) CreateCustomer(ctx context.Context, customer model.Customer) error {
	params := convertModelCustomerToCreateParams(customer)
	err := s.queries.CreateCustomer(ctx, params)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.ErrConflict
			}
		}
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) GetCustomerByID(ctx context.Context, id string) (model.Customer, error) {
	customerID, err := uuid.Parse(id)
	if err != nil {
		return model.Customer{}, errors.ErrBadRequest
	}

	dbCustomer, err := s.queries.GetCustomerByID(ctx, customerID)
	if err == sql.ErrNoRows {
		return model.Customer{}, errors.ErrNotFound
	}
	if err != nil {
		return model.Customer{}, errors.ErrInternalServerError
	}

	return convertDBCustomerToModel(dbCustomer), nil
}

func (s *Storage) GetCustomerByEmail(ctx context.Context, email string) (model.Customer, error) {
	dbCustomer, err := s.queries.GetCustomerByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return model.Customer{}, errors.ErrNotFound
	}
	if err != nil {
		return model.Customer{}, errors.ErrInternalServerError
	}

	return convertDBCustomerToModel(dbCustomer), nil
}

func (s *Storage) ListCustomers(ctx context.Context, limit, offset int) ([]model.Customer, error) {
	params := db.ListCustomersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	dbCustomers, err := s.queries.ListCustomers(ctx, params)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	customers := make([]model.Customer, 0, len(dbCustomers))
	for _, dbCustomer := range dbCustomers {
		customers = append(customers, convertDBCustomerToModel(dbCustomer))
	}

	return customers, nil
}

func (s *Storage) UpdateCustomer(ctx context.Context, customer model.Customer) error {
	_, err := s.queries.GetCustomerByID(ctx, customer.ID)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	if err != nil {
		return errors.ErrInternalServerError
	}

	params := convertModelCustomerToUpdateParams(customer)
	err = s.queries.UpdateCustomer(ctx, params)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.ErrConflict
			}
		}
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) DeleteCustomer(ctx context.Context, id string) error {
	customerID, err := uuid.Parse(id)
	if err != nil {
		return errors.ErrBadRequest
	}

	_, err = s.queries.GetCustomerByID(ctx, customerID)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	if err != nil {
		return errors.ErrInternalServerError
	}

	err = s.queries.DeleteCustomer(ctx, customerID)
	if err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) CreateVendor(ctx context.Context, vendor model.Vendor) error {
	params := convertModelVendorToCreateParams(vendor)
	err := s.queries.CreateVendor(ctx, params)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.ErrConflict
			}
		}
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) GetVendorByID(ctx context.Context, id string) (model.Vendor, error) {
	vendorID, err := uuid.Parse(id)
	if err != nil {
		return model.Vendor{}, errors.ErrBadRequest
	}

	dbVendor, err := s.queries.GetVendorByID(ctx, vendorID)
	if err == sql.ErrNoRows {
		return model.Vendor{}, errors.ErrNotFound
	}
	if err != nil {
		return model.Vendor{}, errors.ErrInternalServerError
	}

	return convertDBVendorToModel(dbVendor), nil
}

func (s *Storage) GetVendorByEmail(ctx context.Context, email string) (model.Vendor, error) {
	dbVendor, err := s.queries.GetVendorByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return model.Vendor{}, errors.ErrNotFound
	}
	if err != nil {
		return model.Vendor{}, errors.ErrInternalServerError
	}

	return convertDBVendorToModel(dbVendor), nil
}

func (s *Storage) ListVendors(ctx context.Context, limit, offset int) ([]model.Vendor, error) {
	params := db.ListVendorsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	dbVendors, err := s.queries.ListVendors(ctx, params)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	vendors := make([]model.Vendor, 0, len(dbVendors))
	for _, dbVendor := range dbVendors {
		vendors = append(vendors, convertDBVendorToModel(dbVendor))
	}

	return vendors, nil
}

func (s *Storage) UpdateVendor(ctx context.Context, vendor model.Vendor) error {
	_, err := s.queries.GetVendorByID(ctx, vendor.ID)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	if err != nil {
		return errors.ErrInternalServerError
	}

	params := convertModelVendorToUpdateParams(vendor)
	err = s.queries.UpdateVendor(ctx, params)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.ErrConflict
			}
		}
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) DeleteVendor(ctx context.Context, id string) error {
	vendorID, err := uuid.Parse(id)
	if err != nil {
		return errors.ErrBadRequest
	}

	_, err = s.queries.GetVendorByID(ctx, vendorID)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	if err != nil {
		return errors.ErrInternalServerError
	}

	err = s.queries.DeleteVendor(ctx, vendorID)
	if err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}
