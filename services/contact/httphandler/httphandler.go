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

// ListCustomers godoc
// @Summary      List customers
// @Description  Get a paginated list of customers
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        limit query int false "Limit" default(10)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} response.SuccessResponse{data=[]model.Customer}
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /customers [get]
// @Security     BearerAuth
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

// GetCustomer godoc
// @Summary      Get customer by ID
// @Description  Get customer details by ID
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        id path string true "Customer ID"
// @Success      200 {object} response.SuccessResponse{data=model.Customer}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /customers/{id} [get]
// @Security     BearerAuth
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

// CreateCustomer godoc
// @Summary      Create customer
// @Description  Create a new customer
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        request body model.CreateCustomerRequest true "Customer creation request"
// @Success      201 {object} response.SuccessResponse{data=model.Customer}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      409 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /customers [post]
// @Security     BearerAuth
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

	customer, err := h.service.CreateCustomer(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to create customer", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, "Customer created successfully", customer, nil)
}

// UpdateCustomer godoc
// @Summary      Update customer
// @Description  Update an existing customer
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        id path string true "Customer ID"
// @Param        request body model.UpdateCustomerRequest true "Customer update request"
// @Success      200 {object} response.SuccessResponse{data=model.Customer}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /customers/{id} [put]
// @Security     BearerAuth
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

	customer, err := h.service.UpdateCustomer(ctx, id, req)
	if err != nil {
		h.logger.Error(ctx, "failed to update customer", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Customer updated successfully", customer, nil)
}

// DeleteCustomer godoc
// @Summary      Delete customer
// @Description  Delete a customer by ID
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        id path string true "Customer ID"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.SimpleErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /customers/{id} [delete]
// @Security     BearerAuth
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

// ListVendors godoc
// @Summary      List vendors
// @Description  Get a paginated list of vendors
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        limit query int false "Limit" default(10)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} response.SuccessResponse{data=[]model.Vendor}
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /vendors [get]
// @Security     BearerAuth
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

// GetVendor godoc
// @Summary      Get vendor by ID
// @Description  Get vendor details by ID
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        id path string true "Vendor ID"
// @Success      200 {object} response.SuccessResponse{data=model.Vendor}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /vendors/{id} [get]
// @Security     BearerAuth
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

// CreateVendor godoc
// @Summary      Create vendor
// @Description  Create a new vendor
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        request body model.CreateVendorRequest true "Vendor creation request"
// @Success      201 {object} response.SuccessResponse{data=model.Vendor}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      409 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /vendors [post]
// @Security     BearerAuth
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

	vendor, err := h.service.CreateVendor(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to create vendor", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, "Vendor created successfully", vendor, nil)
}

// UpdateVendor godoc
// @Summary      Update vendor
// @Description  Update an existing vendor
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        id path string true "Vendor ID"
// @Param        request body model.UpdateVendorRequest true "Vendor update request"
// @Success      200 {object} response.SuccessResponse{data=model.Vendor}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /vendors/{id} [put]
// @Security     BearerAuth
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

	vendor, err := h.service.UpdateVendor(ctx, id, req)
	if err != nil {
		h.logger.Error(ctx, "failed to update vendor", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Vendor updated successfully", vendor, nil)
}

// DeleteVendor godoc
// @Summary      Delete vendor
// @Description  Delete a vendor by ID
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        id path string true "Vendor ID"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.SimpleErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /vendors/{id} [delete]
// @Security     BearerAuth
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
