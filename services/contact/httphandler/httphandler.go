package httphandler

import (
	"encoding/json"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/pagination"
	"microservice-challenge/package/response"
	"microservice-challenge/services/contact/model"
	"microservice-challenge/services/contact/usecase/contact"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	usecase *contact.Usecase
	logger  log.Logger
}

func NewHandler(usecase *contact.Usecase, logger log.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset := pagination.GetLimitOffset(r)

	customers, err := h.usecase.ListCustomers(ctx, limit, offset)
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

	customer, err := h.usecase.GetCustomerByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to get customer", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Customer retrieved successfully", customer, nil)
}

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	customer, err := h.usecase.CreateCustomer(ctx, req)
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
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	customer, err := h.usecase.UpdateCustomer(ctx, id, req)
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

	if err := h.usecase.DeleteCustomer(ctx, id); err != nil {
		h.logger.Error(ctx, "failed to delete customer", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Customer deleted successfully", nil, nil)
}

func (h *Handler) ListVendors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset := pagination.GetLimitOffset(r)

	vendors, err := h.usecase.ListVendors(ctx, limit, offset)
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

	vendor, err := h.usecase.GetVendorByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to get vendor", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Vendor retrieved successfully", vendor, nil)
}

func (h *Handler) CreateVendor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.CreateVendorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	vendor, err := h.usecase.CreateVendor(ctx, req)
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
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.UpdateVendorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	vendor, err := h.usecase.UpdateVendor(ctx, id, req)
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

	if err := h.usecase.DeleteVendor(ctx, id); err != nil {
		h.logger.Error(ctx, "failed to delete vendor", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Vendor deleted successfully", nil, nil)
}
