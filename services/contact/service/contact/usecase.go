package service

import (
	"context"
	"errors"
	pkgerrors "microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/nats"
	"microservice-challenge/services/contact/model"
	"microservice-challenge/services/contact/storage"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	storage    storage.Storage
	natsClient *nats.Client
	logger     log.Logger
}

func NewService(storage storage.Storage, natsClient *nats.Client, logger log.Logger) *Service {
	return &Service{
		storage:    storage,
		natsClient: natsClient,
		logger:     logger,
	}
}

func (s *Service) CreateCustomer(ctx context.Context, req model.CreateCustomerRequest) (model.Customer, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	_, err := s.storage.GetCustomerByEmail(ctx, email)
	if err == nil {
		return model.Customer{}, pkgerrors.ErrConflict
	}
	if !errors.Is(err, pkgerrors.ErrNotFound) {
		return model.Customer{}, err
	}

	customer := model.Customer{
		ID:        uuid.New(),
		Name:      strings.TrimSpace(req.Name),
		Email:     email,
		Phone:     strings.TrimSpace(req.Phone),
		Address:   strings.TrimSpace(req.Address),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.storage.CreateCustomer(ctx, customer); err != nil {
		return model.Customer{}, err
	}

	event := map[string]interface{}{
		"event_type":  "contact.customer.created",
		"customer_id": customer.ID.String(),
		"name":        customer.Name,
		"email":       customer.Email,
		"timestamp":   time.Now().Format(time.RFC3339),
	}
	if err := s.natsClient.Publish("contact.customer.created", event); err != nil {
		s.logger.Error(ctx, "failed to publish customer created event", zap.Error(err))
	}

	return customer, nil
}

func (s *Service) GetCustomerByID(ctx context.Context, id string) (model.Customer, error) {
	return s.storage.GetCustomerByID(ctx, id)
}

func (s *Service) ListCustomers(ctx context.Context, limit, offset int) ([]model.Customer, error) {
	return s.storage.ListCustomers(ctx, limit, offset)
}

func (s *Service) UpdateCustomer(ctx context.Context, id string, req model.UpdateCustomerRequest) (model.Customer, error) {
	customer, err := s.storage.GetCustomerByID(ctx, id)
	if err != nil {
		return model.Customer{}, err
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email != customer.Email {
		_, err := s.storage.GetCustomerByEmail(ctx, email)
		if err == nil {
			return model.Customer{}, pkgerrors.ErrConflict
		}
		if !errors.Is(err, pkgerrors.ErrNotFound) {
			return model.Customer{}, err
		}
	}

	customer.Name = strings.TrimSpace(req.Name)
	customer.Email = email
	customer.Phone = strings.TrimSpace(req.Phone)
	customer.Address = strings.TrimSpace(req.Address)
	customer.UpdatedAt = time.Now()

	if err := s.storage.UpdateCustomer(ctx, customer); err != nil {
		return model.Customer{}, err
	}

	event := map[string]interface{}{
		"event_type":  "contact.customer.updated",
		"customer_id": customer.ID.String(),
		"name":        customer.Name,
		"email":       customer.Email,
		"timestamp":   time.Now().Format(time.RFC3339),
	}
	if err := s.natsClient.Publish("contact.customer.updated", event); err != nil {
		s.logger.Error(ctx, "failed to publish customer updated event", zap.Error(err))
	}

	return customer, nil
}

func (s *Service) DeleteCustomer(ctx context.Context, id string) error {
	return s.storage.DeleteCustomer(ctx, id)
}

func (s *Service) CreateVendor(ctx context.Context, req model.CreateVendorRequest) (model.Vendor, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	_, err := s.storage.GetVendorByEmail(ctx, email)
	if err == nil {
		return model.Vendor{}, pkgerrors.ErrConflict
	}
	if !errors.Is(err, pkgerrors.ErrNotFound) {
		return model.Vendor{}, err
	}

	vendor := model.Vendor{
		ID:        uuid.New(),
		Name:      strings.TrimSpace(req.Name),
		Email:     email,
		Phone:     strings.TrimSpace(req.Phone),
		Address:   strings.TrimSpace(req.Address),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.storage.CreateVendor(ctx, vendor); err != nil {
		return model.Vendor{}, err
	}

	event := map[string]interface{}{
		"event_type": "contact.vendor.created",
		"vendor_id":  vendor.ID.String(),
		"name":       vendor.Name,
		"email":      vendor.Email,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	if err := s.natsClient.Publish("contact.vendor.created", event); err != nil {
		s.logger.Error(ctx, "failed to publish vendor created event", zap.Error(err))
	}

	return vendor, nil
}

func (s *Service) GetVendorByID(ctx context.Context, id string) (model.Vendor, error) {
	return s.storage.GetVendorByID(ctx, id)
}

func (s *Service) ListVendors(ctx context.Context, limit, offset int) ([]model.Vendor, error) {
	return s.storage.ListVendors(ctx, limit, offset)
}

func (s *Service) UpdateVendor(ctx context.Context, id string, req model.UpdateVendorRequest) (model.Vendor, error) {
	vendor, err := s.storage.GetVendorByID(ctx, id)
	if err != nil {
		return model.Vendor{}, err
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email != vendor.Email {
		_, err := s.storage.GetVendorByEmail(ctx, email)
		if err == nil {
			return model.Vendor{}, pkgerrors.ErrConflict
		}
		if !errors.Is(err, pkgerrors.ErrNotFound) {
			return model.Vendor{}, err
		}
	}

	vendor.Name = strings.TrimSpace(req.Name)
	vendor.Email = email
	vendor.Phone = strings.TrimSpace(req.Phone)
	vendor.Address = strings.TrimSpace(req.Address)
	vendor.UpdatedAt = time.Now()

	if err := s.storage.UpdateVendor(ctx, vendor); err != nil {
		return model.Vendor{}, err
	}

	event := map[string]interface{}{
		"event_type": "contact.vendor.updated",
		"vendor_id":  vendor.ID.String(),
		"name":       vendor.Name,
		"email":      vendor.Email,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	if err := s.natsClient.Publish("contact.vendor.updated", event); err != nil {
		s.logger.Error(ctx, "failed to publish vendor updated event", zap.Error(err))
	}

	return vendor, nil
}

func (s *Service) DeleteVendor(ctx context.Context, id string) error {
	return s.storage.DeleteVendor(ctx, id)
}
