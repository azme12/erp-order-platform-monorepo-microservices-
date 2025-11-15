package postgresql

import (
	"context"
	"database/sql"
	"microservice-challenge/package/errors"
	"microservice-challenge/services/contact/model"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) DB() *sql.DB {
	return s.db
}

func (s *Storage) CreateCustomer(ctx context.Context, customer model.Customer) error {
	query := `
		INSERT INTO customers (id, name, email, phone, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := s.db.ExecContext(ctx, query,
		customer.ID,
		customer.Name,
		strings.ToLower(strings.TrimSpace(customer.Email)),
		customer.Phone,
		customer.Address,
		customer.CreatedAt,
		customer.UpdatedAt,
	)

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
	var customer model.Customer

	customerID, err := uuid.Parse(id)
	if err != nil {
		return customer, errors.ErrBadRequest
	}

	query := `
		SELECT id, name, email, phone, address, created_at, updated_at
		FROM customers
		WHERE id = $1
	`

	err = s.db.QueryRowContext(ctx, query, customerID).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return customer, errors.ErrNotFound
	}

	if err != nil {
		return customer, errors.ErrInternalServerError
	}

	return customer, nil
}

func (s *Storage) GetCustomerByEmail(ctx context.Context, email string) (model.Customer, error) {
	var customer model.Customer

	query := `
		SELECT id, name, email, phone, address, created_at, updated_at
		FROM customers
		WHERE email = $1
	`

	err := s.db.QueryRowContext(ctx, query, strings.ToLower(strings.TrimSpace(email))).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return customer, errors.ErrNotFound
	}

	if err != nil {
		return customer, errors.ErrInternalServerError
	}

	return customer, nil
}

func (s *Storage) ListCustomers(ctx context.Context, limit, offset int) ([]model.Customer, error) {
	query := `
		SELECT id, name, email, phone, address, created_at, updated_at
		FROM customers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	defer rows.Close()

	customers := make([]model.Customer, 0, limit)
	for rows.Next() {
		var customer model.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.Address,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, errors.ErrInternalServerError
		}
		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.ErrInternalServerError
	}

	return customers, nil
}

func (s *Storage) UpdateCustomer(ctx context.Context, customer model.Customer) error {
	query := `
		UPDATE customers
		SET name = $2, email = $3, phone = $4, address = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := s.db.ExecContext(ctx, query,
		customer.ID,
		customer.Name,
		strings.ToLower(strings.TrimSpace(customer.Email)),
		customer.Phone,
		customer.Address,
		customer.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.ErrConflict
			}
		}
		return errors.ErrInternalServerError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.ErrInternalServerError
	}

	if rowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (s *Storage) DeleteCustomer(ctx context.Context, id string) error {
	customerID, err := uuid.Parse(id)
	if err != nil {
		return errors.ErrBadRequest
	}

	query := `DELETE FROM customers WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, customerID)
	if err != nil {
		return errors.ErrInternalServerError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.ErrInternalServerError
	}

	if rowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (s *Storage) CreateVendor(ctx context.Context, vendor model.Vendor) error {
	query := `
		INSERT INTO vendors (id, name, email, phone, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := s.db.ExecContext(ctx, query,
		vendor.ID,
		vendor.Name,
		strings.ToLower(strings.TrimSpace(vendor.Email)),
		vendor.Phone,
		vendor.Address,
		vendor.CreatedAt,
		vendor.UpdatedAt,
	)

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
	var vendor model.Vendor

	vendorID, err := uuid.Parse(id)
	if err != nil {
		return vendor, errors.ErrBadRequest
	}

	query := `
		SELECT id, name, email, phone, address, created_at, updated_at
		FROM vendors
		WHERE id = $1
	`

	err = s.db.QueryRowContext(ctx, query, vendorID).Scan(
		&vendor.ID,
		&vendor.Name,
		&vendor.Email,
		&vendor.Phone,
		&vendor.Address,
		&vendor.CreatedAt,
		&vendor.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return vendor, errors.ErrNotFound
	}

	if err != nil {
		return vendor, errors.ErrInternalServerError
	}

	return vendor, nil
}

func (s *Storage) GetVendorByEmail(ctx context.Context, email string) (model.Vendor, error) {
	var vendor model.Vendor

	query := `
		SELECT id, name, email, phone, address, created_at, updated_at
		FROM vendors
		WHERE email = $1
	`

	err := s.db.QueryRowContext(ctx, query, strings.ToLower(strings.TrimSpace(email))).Scan(
		&vendor.ID,
		&vendor.Name,
		&vendor.Email,
		&vendor.Phone,
		&vendor.Address,
		&vendor.CreatedAt,
		&vendor.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return vendor, errors.ErrNotFound
	}

	if err != nil {
		return vendor, errors.ErrInternalServerError
	}

	return vendor, nil
}

func (s *Storage) ListVendors(ctx context.Context, limit, offset int) ([]model.Vendor, error) {
	query := `
		SELECT id, name, email, phone, address, created_at, updated_at
		FROM vendors
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	defer rows.Close()

	vendors := make([]model.Vendor, 0, limit)
	for rows.Next() {
		var vendor model.Vendor
		err := rows.Scan(
			&vendor.ID,
			&vendor.Name,
			&vendor.Email,
			&vendor.Phone,
			&vendor.Address,
			&vendor.CreatedAt,
			&vendor.UpdatedAt,
		)
		if err != nil {
			return nil, errors.ErrInternalServerError
		}
		vendors = append(vendors, vendor)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.ErrInternalServerError
	}

	return vendors, nil
}

func (s *Storage) UpdateVendor(ctx context.Context, vendor model.Vendor) error {
	query := `
		UPDATE vendors
		SET name = $2, email = $3, phone = $4, address = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := s.db.ExecContext(ctx, query,
		vendor.ID,
		vendor.Name,
		strings.ToLower(strings.TrimSpace(vendor.Email)),
		vendor.Phone,
		vendor.Address,
		vendor.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.ErrConflict
			}
		}
		return errors.ErrInternalServerError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.ErrInternalServerError
	}

	if rowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (s *Storage) DeleteVendor(ctx context.Context, id string) error {
	vendorID, err := uuid.Parse(id)
	if err != nil {
		return errors.ErrBadRequest
	}

	query := `DELETE FROM vendors WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, vendorID)
	if err != nil {
		return errors.ErrInternalServerError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.ErrInternalServerError
	}

	if rowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}
