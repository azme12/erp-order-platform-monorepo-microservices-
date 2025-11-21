package httphandler

import (
	"encoding/json"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/pagination"
	"microservice-challenge/package/response"
	"microservice-challenge/services/contact/model"
	contactservice "microservice-challenge/services/contact/service/contact"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	maxRequestBodySize = 1 << 20
)

type Handler struct {
	service *contactservice.Service
	logger  log.Logger
}

func NewHandler(service *contactservice.Service, logger log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) parseAndValidateRequest(w http.ResponseWriter, r *http.Request, req interface{ Validate() error }) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return err
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return err
	}

	return nil
}

func (h *Handler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset := pagination.GetLimitOffset(r)

	customers, err := h.service.ListCustomers(ctx, limit, offset)
	if err != nil {
		h.logger.Error(ctx, "failed to list customers", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Customers retrieved successfully", customers, nil)
}

func (h *Handler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	customer, err := h.service.GetCustomerByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to get customer", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Customer retrieved successfully", customer, nil)
}

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.CreateCustomerRequest
	if err := h.parseAndValidateRequest(w, r, &req); err != nil {
		return
	}

	customer, err := h.service.CreateCustomer(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to create customer", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, "Customer created successfully", customer, nil)
}

func (h *Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var req model.UpdateCustomerRequest
	if err := h.parseAndValidateRequest(w, r, &req); err != nil {
		return
	}

	customer, err := h.service.UpdateCustomer(ctx, id, req)
	if err != nil {
		h.logger.Error(ctx, "failed to update customer", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Customer updated successfully", customer, nil)
}

func (h *Handler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteCustomer(ctx, id); err != nil {
		h.logger.Error(ctx, "failed to delete customer", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Customer deleted successfully", nil, nil)
}

func (h *Handler) ListVendors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset := pagination.GetLimitOffset(r)

	vendors, err := h.service.ListVendors(ctx, limit, offset)
	if err != nil {
		h.logger.Error(ctx, "failed to list vendors", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Vendors retrieved successfully", vendors, nil)
}

func (h *Handler) GetVendor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	vendor, err := h.service.GetVendorByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to get vendor", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Vendor retrieved successfully", vendor, nil)
}

func (h *Handler) CreateVendor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.CreateVendorRequest
	if err := h.parseAndValidateRequest(w, r, &req); err != nil {
		return
	}

	vendor, err := h.service.CreateVendor(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to create vendor", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, "Vendor created successfully", vendor, nil)
}

func (h *Handler) UpdateVendor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var req model.UpdateVendorRequest
	if err := h.parseAndValidateRequest(w, r, &req); err != nil {
		return
	}

	vendor, err := h.service.UpdateVendor(ctx, id, req)
	if err != nil {
		h.logger.Error(ctx, "failed to update vendor", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Vendor updated successfully", vendor, nil)
}

func (h *Handler) DeleteVendor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteVendor(ctx, id); err != nil {
		h.logger.Error(ctx, "failed to delete vendor", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Vendor deleted successfully", nil, nil)
}
